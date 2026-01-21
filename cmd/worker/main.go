package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"pointservice/internal/infra/aync/rabbitmq"
	"pointservice/internal/infra/database/mysql"
	"pointservice/internal/infra/repository"
	"pointservice/internal/presentation"
	"pointservice/internal/usecase"
	"sync"
	"syscall"
)

const (
	defaultEnvironment = "dev"
)

// main 依存関係の呼び出し, 環境変数, アプリケーションの起動, クリーンシャットダウン
func main() {
	var environment = defaultEnvironment
	if v := os.Getenv("ENVIRONMENT"); v != "" {
		environment = v
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	db, closeDB := mysql.ConnectDB()
	defer closeDB()
	pointRepo := repository.NewPointSQL(db)
	reservationRepo := repository.NewReservationSQL(db)
	addReservationUseCase := usecase.NewAddReservationPointUseCase(pointRepo, reservationRepo)
	pointWorkerHandler := presentation.NewPointWorkerHandler(addReservationUseCase)

	log.Println("Starting point grant worker...")
	conn := rabbitmq.NewConnection(true, environment)
	defer conn.Conn.Close()
	defer conn.Ch.Close()
	queue, err := rabbitmq.NewQueueDeclare(conn.Ch)
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := rabbitmq.NewConsume(conn.Ch, queue)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Worker waiting for messages...")
	var wg sync.WaitGroup
	go func() {
		defer wg.Done()
		pointWorkerHandler.PointReserveWorker(ctx, msgs)
	}()
	<-ctx.Done()
	log.Println("Stopping point grant worker...")
	wg.Wait()
	log.Println("Worker gracefully stopped")
}
