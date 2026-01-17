package infra

import (
	"pointservice/internal/presentation"

	"github.com/labstack/echo/v4"
)

type Router struct {
	server *echo.Echo
}

func NewRouter(server *echo.Echo) Router {
	return Router{
		server: server,
	}
}

func (r *Router) Exec(h *presentation.PointHandler) {
	e := r.server
	e.PUT("/point/add", h.PointAdd)
	e.PUT("/point/sub", h.PointSub)
	e.GET("/point/confirm", h.PointConfirm)
	e.POST("/point/reserve", h.PointReserve)
	e.GET("/", h.HealthCheck)
	e.GET("/health", h.HealthCheck)
	e.GET("/health/", h.HealthCheck)
}
