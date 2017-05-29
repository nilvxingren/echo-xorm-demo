package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nilvxingren/echoxormdemo/app"
	"github.com/nilvxingren/echoxormdemo/ctx"
)

var (
	configFlag = flag.String("config",
		"./resource/config.toml",
		"-config=\"path-to-your-config-file\" ")
)

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	// parse flags
	flag.Parse()

	var (
		err error
		a   *app.Application
	)

	flags := &ctx.Flags{
		CfgFileName: *configFlag,
	}

	// create application
	a, err = app.New(flags)
	if err != nil {
		log.Fatal("error ", os.Args[0]+" initialization error: "+err.Error())
		os.Exit(1)
	}
	// setup OS-signal catchers
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() { // start OS-signal catching route
		for sig := range signalChannel {
			if a.C.Orm != nil {
				err = a.C.Orm.Close()
				if err != nil {
					a.C.Logger.Error("appcontrol", os.Args[0]+" db closing error on "+sig.String())
				}
			}
			if a.C.Logger != nil {
				a.C.Logger.Info("appcontrol", os.Args[0]+" graceful shutdown on "+sig.String())
				a.C.Logger.Close()
			}
			os.Exit(1)
		}
	}()

	// run application server
	if a.C.Logger == nil {
		log.Fatal("error ", os.Args[0]+" startup error: logger not initialized ")
		os.Exit(1)
	}
	a.C.Logger.Info("appcontrol", "started on localhost:"+a.C.Config.Port)
	a.Run()
}
