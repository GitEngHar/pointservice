package main

import (
	"pointservice/internal/adapter/repository"
	"pointservice/internal/infra"
	"pointservice/internal/infra/database/mysql"
	"pointservice/internal/presentation"
)

func main() {
	db, closer := mysql.ConnectDB()
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
