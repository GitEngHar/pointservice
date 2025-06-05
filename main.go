package main

import "pointservice/infra"

func main() {
	var app = infra.NewConfig().
		WebServer()
	app.Start()
}
