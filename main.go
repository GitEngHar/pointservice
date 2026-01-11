package main

import (
	"pointservice/internal/infra"
	"pointservice/internal/infra/aync/mq"
	"pointservice/internal/infra/database/mysql"
	"pointservice/internal/infra/repository"
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
	producer, closer := mq.ConnectProducer()
	defer func() {
		if closer == nil {
			return
		}
		_ = closer()
	}()
	tallyProducer := mq.NewRabbitProducer(producer)
	handler := presentation.NewPointHandler(db, repo, tallyProducer)

	var app = infra.NewConfig().
		WebServer()
	app.Start(handler)
}
