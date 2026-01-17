package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"

	"pointservice/internal/domain"
	"pointservice/internal/infra/aync/mq"
	"pointservice/internal/infra/database/mysql"
	"pointservice/internal/infra/repository"
)

const (
	reservationQueueName = "reservationQueue"
	rabbitURI            = "amqp://guest:guest@rabbitmq:5672/"
)

func main() {
	log.Println("Starting point grant worker...")

	// Connect to DB
	db, dbCloser := mysql.ConnectDB()
	defer func() {
		_ = dbCloser()
	}()

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare queue
	queue, err := ch.QueueDeclare( // キュー（ポスト）が存在するか確認、なければ作る。
		reservationQueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume( // このキューから予約メッセージを受け取る。
		queue.Name,
		"",    // consumer tag
		false, // auto-ack (manual ack for reliability)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	pointRepo := repository.NewPointSQL(db)
	reservationRepo := repository.NewReservationSQL(db)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Worker waiting for messages...")

	go func() {
		for msg := range msgs {
			if err := processMessage(ctx, msg, pointRepo, reservationRepo); err != nil {
				log.Printf("Error processing message: %v", err)

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

	sig := <-sigCh
	log.Printf("Received signal %v, shutting down...", sig)
}

// 予約メッセージを処理する
func processMessage(
	ctx context.Context,
	msg amqp.Delivery,
	pointRepo repository.PointRepository,
	reservationRepo repository.ReservationRepository,
) error {
	var reservation mq.ReservationMessage
	if err := json.Unmarshal(msg.Body, &reservation); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return err
	}

	log.Printf("Processing reservation: %s (user: %s, amount: %d)",
		reservation.ReservationID, reservation.UserID, reservation.PointAmount)

	// Idempotent point grant
	added, err := pointRepo.AddPointIdempotent(
		ctx,
		reservation.IdempotencyKey,
		reservation.UserID,
		reservation.PointAmount,
	)
	if err != nil {
		log.Printf("Failed to add point for reservation %s: %v", reservation.ReservationID, err)
		// Update reservation status to FAILED
		if updateErr := reservationRepo.UpdateStatus(ctx, reservation.ReservationID, domain.StatusFailed); updateErr != nil {
			log.Printf("Failed to update reservation status to FAILED: %v", updateErr)
		}
		return err
	}

	if !added {

		// すでに処理済みだった場合のログ。
		log.Printf("Reservation %s already processed (idempotency key: %s)",
			reservation.ReservationID, reservation.IdempotencyKey)
	} else {
		log.Printf("Successfully granted %d points to user %s",
			reservation.PointAmount, reservation.UserID)
	}

	// Update reservation status to DONE
	if err := reservationRepo.UpdateStatus(ctx, reservation.ReservationID, domain.StatusDone); err != nil {
		log.Printf("Failed to update reservation status to DONE: %v", err)
		return err
	}

	log.Printf("Reservation %s completed", reservation.ReservationID)
	return nil
}
