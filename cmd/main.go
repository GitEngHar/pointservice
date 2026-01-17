package main

import (
	"pointservice/internal/infra"
	"pointservice/internal/infra/database/mysql"
	"pointservice/internal/infra/repository"
	"pointservice/internal/presentation"
)

func main() {
	db, closer := mysql.ConnectDB()
	defer func() {
		_ = closer()
	}()
	repo := repository.NewPointSQL(db)
	reservationRepo := repository.NewReservationSQL(db)
	defer func() {
		if closer != nil {
			_ = closer()
		}
	}()
	handler := presentation.NewPointHandler(db, repo, reservationRepo)

	var app = infra.NewConfig().
		WebServer()
	app.Start(handler)
}
