package infra

import "pointservice/adapter/repository"
import "github.com/labstack/echo/v4"

type config struct {
	appName string
	server  *echo.Echo
	router  Router
	dbSQL   repository.PointRepository
}

func NewConfig() *config {
	return &config{}
}

func (c *config) WebServer() *config {
	c.server = echo.New()
	c.router = NewRouter(c.server)
	return c
}

func (c *config) Start() {
	e := c.server
	c.router.exec()
	e.Logger.Fatal(e.Start(":1323"))
}
