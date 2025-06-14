package infra

import (
	"pointservice/adapter/repository"
	"pointservice/infra/log"
	"pointservice/presentation"
)
import "github.com/labstack/echo/v4"

type Config struct {
	appName string
	server  *echo.Echo
	router  Router
	log     log.RequestLogger
	dbSQL   repository.PointRepository
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) WebServer() *Config {
	c.server = echo.New()
	c.router = NewRouter(c.server)
	c.log = log.NewRequestLogger(c.server)
	return c
}

func (c *Config) Start(h *presentation.PointHandler) {
	e := c.server
	c.router.Exec(h)
	c.log.Exec()
	e.Logger.Fatal(e.Start(":1323"))
}
