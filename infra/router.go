package infra

import (
	"github.com/labstack/echo/v4"
	"pointservice/presentation"
)

type Router struct {
	server *echo.Echo
}

func NewRouter(server *echo.Echo) Router {
	return Router{
		server: server,
	}
}

func (r *Router) exec(h *presentation.PointHandler) {
	e := r.server
	e.PUT("/point/add", h.PointAdd)
	e.PUT("/point/sub", h.PointSub)
}
