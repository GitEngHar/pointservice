package main

import (
	"pointservice/adapter/repository"
	"pointservice/infra"
	"pointservice/infra/database"
	"pointservice/presentation"
)

func main() {
	db, closer := database.ConnectDB()
	defer func() {
		if err := closer(); err != nil {
			panic(err)
		}
	}()
	repo := repository.NewPointSQL(db)
	handler := presentation.NewPointHandler(db, repo)

	var app = infra.NewConfig().
		WebServer()
	app.Start(handler)
}
