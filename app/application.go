package app

import (
	"github.com/nilvxingren/echoxormdemo/ctx"
	"github.com/nilvxingren/echoxormdemo/server"
	"net/http"
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"strconv"
	"os"
	"errors"

	"github.com/nilvxingren/echoxormdemo/logger"
	"github.com/go-xorm/xorm"
	"github.com/nilvxingren/echoxormdemo/server/users"
)

// Application define a mode of running app
type Application struct {
	C *ctx.Context
}

// New constructor
func New(flags *ctx.Flags) (*Application, error) {
	app := new(Application)
	app.C = new(ctx.Context)
	// read config file
	err := app.initConfigFromFile(flags.CfgFileName)
	if err != nil {
		return nil, err
	}

	// init Logger
	err = app.initLogger()
	if err != nil {
		return nil, err
	}

	// init Orm
	err = app.initOrm()
	return app, err
}

// Run starts application
func (a *Application) Run() {
	srv := server.New(a.C)
	srv.Run()
}


// readConfig reads configuration file into application Config structure and inits in-memory token storage
func (a *Application) initConfigFromFile(cfgFileName string) error {
	// read config
	tomlData, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		return errors.New("Configuration file read error: " + cfgFileName + "\nError:" + err.Error())
	}
	_, err = toml.Decode(string(tomlData[:]), &a.C.Config)
	if err != nil {
		return errors.New("Configuration file decoding error: " + cfgFileName + "\nError:" + err.Error())
	}
	// init Logging data
	if len(a.C.Config.Logging.ID) == 0 {
		a.C.Config.Logging.ID = strconv.Itoa(os.Getpid())
	}
	if len(a.C.Config.Logging.LogTag) == 0 {
		a.C.Config.Logging.LogTag = os.Args[0]
	}
	return nil
}

// setupLogger sets apllication Logger up according to configuration settings
func (a *Application) initLogger() error {
	if a.C.Config.Logging.LogMode == "nil" || a.C.Config.Logging.LogMode == "null" {
		a.C.Logger = logger.NewNilLogger()
		return nil
	}
	a.C.Logger = logger.NewStdLogger(a.C.Config.Logging.ID, a.C.Config.Logging.LogTag)
	return nil
}

// init database
func (a *Application) initOrm() error {
	var err error
	// open database
	a.C.Orm, err = xorm.NewEngine(a.C.Config.Database.Db, a.C.Config.Database.Dsn)
	if err != nil {
		return err
	}
	// turn on logs
	ormLogger := logger.NewOrmLogger(a.C.Logger)
	a.C.Orm.SetLogger(ormLogger)
	a.C.Orm.ShowSQL(true)
	// migrate
	//err = a.migrateDb()
	//if err != nil {
	//	return err
	//}
	//// init data
	//err = a.initDbData()
	return err
}

// migrate database
func (a *Application) migrateDb() error {
	var err error
	// migrate tables
	err = a.C.Orm.Sync(new(users.User))
	return err
}

// initDbData installs hardcoded data from config
func (a *Application) initDbData() error {
	user := &users.User{Login: "admin", Password: "admin"} // aaaa, backdoor
	status, err := user.Save(a.C.Orm)
	if err == nil {
		return nil
	}
	if status == http.StatusConflict {
		return nil
	}
	return err
}