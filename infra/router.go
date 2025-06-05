package infra

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Router struct {
	server *echo.Echo
}

func NewRouter(server *echo.Echo) Router {
	return Router{
		server: server,
	}
}

func (r Router) pointAdd(c echo.Context) error {
	returnMassage := "pointAdd"
	return c.String(http.StatusOK, returnMassage)
}

func (r Router) pointSub(c echo.Context) error {
	returnMassage := "pointSub"
	return c.String(http.StatusOK, returnMassage)
}

func (r Router) exec() {
	e := r.server
	e.GET("/point/add", r.pointAdd)
	e.GET("/point/sub", r.pointSub)
}
