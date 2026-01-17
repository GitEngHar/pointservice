package main

import (
	"os"
	"pointservice/internal/infra"
	"pointservice/internal/infra/aync/mq"
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
	environment := os.Getenv("ENVIRONMENT")
	producer, closer := mq.ConnectProducer(environment)
	defer func() {
		if closer != nil {
			_ = closer()
		}
	}()
	tallyProducer := mq.NewRabbitProducer(producer)
	handler := presentation.NewPointHandler(db, repo, reservationRepo, tallyProducer)

	var app = infra.NewConfig().
		WebServer()
	app.Start(handler)
}
