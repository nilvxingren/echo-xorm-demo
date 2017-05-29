package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/nilvxingren/echoxormdemo/ctx"
)

// Input represents payload data format
type Input struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Handler is a container for handlers and app data
type Handler struct {
	C *ctx.Context
}

// GetAllUsers is a GET /users handler
func (h *Handler) GetAllUsers(c echo.Context) error {
	users, err := new(User).FindAll(h.C.Orm)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

// GetUser is a GET /users/{id} handler
func (h *Handler) GetUser(c echo.Context) error {
	var (
		user   User
		err    error
		status int
	)

	user.ID, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	status, err = user.Find(h.C.Orm)
	if err != nil {
		return c.String(status, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// CreateUser is a POST /users handler
func (h *Handler) CreateUser(c echo.Context) error {
	var (
		status int
		err    error
		user   User
		input  Input
	)

	if err = c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// validate
	if len(input.Login) == 0 {
		return c.String(http.StatusBadRequest, "login not recognized")
	}
	if len(input.Password) == 0 {
		return c.String(http.StatusBadRequest, "password not recognized")
	}

	// create
	user = User{
		Login:    input.Login,
		Password: input.Password,
	}
	// save
	status, err = user.Save(h.C.Orm)
	if err != nil {
		return c.String(status, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}

// PutUser is a PUT /users/{id} handler
func (h *Handler) PutUser(c echo.Context) error {
	var (
		input  Input
		user   User
		id     uint64
		err    error
		status int
	)
	// parse id
	id, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// parse request body
	if err = c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// construct user
	user = User{
		ID:       id,
		Login:    input.Login,
		Password: input.Password,
	}
	// update
	status, err = user.Update(h.C.Orm)
	if err != nil {
		return c.String(status, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUser is a DELETE /users/{id} handler
func (h *Handler) DeleteUser(c echo.Context) error {
	var (
		id     uint64
		status int
		err    error
		user   User
	)

	id, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user.ID = id
	// delete
	status, err = user.Delete(h.C.Orm)
	if err != nil {
		return c.String(status, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
