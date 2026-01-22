package main

import (
	"context"
	"log"
	"os"
	"pointservice/internal/infra/aync/rabbitmq"
	"pointservice/internal/infra/database/mysql"
	"pointservice/internal/infra/repository"
	"pointservice/internal/presentation"
	"pointservice/internal/usecase"
)

// スキャン（見回り）をする間隔。ここでは「10秒に1回」。
const (
	defaultEnvironment = "dev"
)

func main() {
	var environment = defaultEnvironment
	if v := os.Getenv("ENVIRONMENT"); v != "" {
		environment = v
	}
	log.Println("Starting reservation scheduler...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, closeDB := mysql.ConnectDB()
	defer closeDB()
	conn := rabbitmq.NewConnection(false, environment)
	defer conn.Conn.Close()
	reservationRepo := repository.NewReservationSQL(db)
	pointRepo := repository.NewPointRepository(db)
	producer := rabbitmq.NewRabbitProducer(conn.Conn)
	pointReservationConfirmUseCase := usecase.NewPointReservationConfirmUseCase(pointRepo, reservationRepo)
	handler := presentation.NewPointSchedulerHandler(pointReservationConfirmUseCase)
	handler.PointReserveConfirmScheduler(ctx, producer)
}
