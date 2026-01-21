package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"pointservice/internal/infra/aync/rabbitmq"
	"pointservice/internal/infra/database/mysql"
	"pointservice/internal/infra/repository"
	"pointservice/internal/usecase"
	"syscall"
)

const (
	reservationQueueName = "reservationQueue"
)

func main() {
	log.Println("Starting point grant worker...")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, closeDB := mysql.ConnectDB()
	defer closeDB()
	pointRepo := repository.NewPointSQL(db)
	reservationRepo := repository.NewReservationSQL(db)
	addReservationUseCase := usecase.NewAddReservationPointUseCase(pointRepo, reservationRepo)
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "dev"
	}

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

	go func() {
		for msg := range msgs {
			if er := addReservationUseCase.Execute(ctx, msg); er != nil {
				log.Printf("Error processing message: %v", er)
				// このメッセージの処理に失敗しました」とキュー（例：RabbitMQ など）に通知する
				if nackErr := msg.Nack(false, true); nackErr != nil {
					log.Printf("Failed to nack message: %v", nackErr)
				}
				continue
			}
			// Ack on success
			if ackErr := msg.Ack(false); ackErr != nil {
				log.Printf("Failed to ack message: %v", ackErr)
			}
		}
	}()

	log.Printf("Received signal %v, shutting down...", sig)
}

// TODO useCaseとhandler
