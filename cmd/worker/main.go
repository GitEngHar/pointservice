package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"pointservice/internal/infra/aync/dto"
	"pointservice/internal/infra/aync/mq"
	"pointservice/internal/usecase"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"

	"pointservice/internal/domain"
	"pointservice/internal/infra/database/mysql"
	"pointservice/internal/infra/repository"
)

const (
	reservationQueueName = "reservationQueue"
)

func main() {

	// TODO infra層へ
	log.Println("Starting point grant worker...")

	db, closeDB := mysql.ConnectDB()
	defer closeDB()

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "dev"
	}
	conn := mq.NewConnection(true, environment)
	defer conn.Conn.Close()
	defer conn.Ch.Close()

	//TODO Infra Consumer
	//queue, err := ch.QueueDeclare( // キュー（ポスト）が存在するか確認、なければ作る。
	//	reservationQueueName,
	//	true,  // durable
	//	false, // delete when unused
	//	false, // exclusive
	//	false, // no-wait
	//	nil,   // arguments
	//)
	//if err != nil {
	//	log.Fatalf("Failed to declare queue: %v", err)
	//}
	//
	//msgs, err := ch.Consume( // このキューから予約メッセージを受け取る。
	//	queue.Name,
	//	"",    // consumer tag
	//	false, // auto-ack (manual ack for reliability)
	//	false, // exclusive
	//	false, // no-local
	//	false, // no-wait
	//	nil,   // arguments
	//)
	//if err != nil {
	//	log.Fatalf("Failed to register consumer: %v", err)
	//}
	//TODO Learn
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pointRepo := repository.NewPointSQL(db)
	reservationRepo := repository.NewReservationSQL(db)

	usecase.NewPointUpsertInterceptor()

	log.Println("Worker waiting for messages...")

	// TODO infra層へ
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

// TODO useCaseとhandler
// 予約メッセージを処理する
func processMessage(
	ctx context.Context,
	msg amqp.Delivery,
	pointRepo repository.PointRepository,
	reservationRepo repository.ReservationRepository,
) error {
	var reservation dto.ReservationMessage
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
