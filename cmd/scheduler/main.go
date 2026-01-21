package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"pointservice/internal/infra/aync/dto"
	"pointservice/internal/infra/aync/rabbitmq"
	"syscall"
	"time"

	"pointservice/internal/domain"
	"pointservice/internal/infra/aync/mq"
	"pointservice/internal/infra/database/mysql"
	"pointservice/internal/infra/repository"
)

// スキャン（見回り）をする間隔。ここでは「10秒に1回」。
const (
	defaultEnvironment = "dev"
	scanInterval       = 10 * time.Second
)

func main() {
	var environment = defaultEnvironment
	if v := os.Getenv("ENVIRONMENT"); v != "" {
		environment = v
	}
	log.Println("Starting reservation scheduler...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Connect to DB
	db, closeDB := mysql.ConnectDB()
	defer closeDB()
	conn := rabbitmq.NewConnection(false, environment)
	defer conn.Conn.Close()

	reservationRepo := repository.NewReservationSQL(db)
	producer := rabbitmq.NewRabbitProducer(conn.Conn)
	// TODO usecase追加
	//TODO handler追加する

	// Setup graceful shutdown

	// OSからの終了合図（Ctrl+Cとか）。
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()

	log.Printf("Scheduler running, scanning every %v\n", scanInterval)

	//TODO handler移行
	for {
		select {
		case <-ticker.C: // 10秒経ったら、中の処理をやる。
			if err := scanAndPublish(ctx, reservationRepo, producer); err != nil {
				log.Printf("Error during scan: %v\n", err)
			}
		case sig := <-sigCh: // 終了合図が来たら、終了する。
			log.Printf("Received signal %v, shutting down...\n", sig)
			return
		}
	}
}

// 予約をスキャンして、メッセージをキューに流す。
func scanAndPublish(ctx context.Context, repo repository.ReservationRepository, producer *rabbitmq.RabbitProducer) error {
	now := time.Now()
	reservations, err := repo.GetPendingReservations(ctx, now) // 実行待ちの予約を取得する。
	if err != nil {
		return fmt.Errorf("failed to get pending reservations: %w", err)
	}

	if len(reservations) == 0 {
		log.Println("No pending reservations found")
		return nil
	}

	log.Printf("Found %d pending reservations\n", len(reservations))

	// 予約を一つずつ取り出して、キューに流す。
	for _, res := range reservations {
		// Update status to PROCESSING first
		if err := repo.UpdateStatus(ctx, res.ID, domain.StatusProcessing); err != nil {
			log.Printf("Failed to update reservation %s to PROCESSING: %v\n", res.ID, err)
			continue
		}

		// Publish to queue
		msg := dto.ReservationMessage{
			ReservationID:  res.ID,
			UserID:         res.UserID,
			PointAmount:    res.PointAmount,
			IdempotencyKey: res.IdempotencyKey,
		}

		// RabbitMQにメッセージを送る（Publish）。
		if err := producer.PublishReservation(ctx, msg); err != nil {
			log.Printf("Failed to publish reservation %s: %v\n", res.ID, err)
			// Revert status back to PENDING on publish failure
			if revertErr := repo.UpdateStatus(ctx, res.ID, domain.StatusPending); revertErr != nil {
				log.Printf("Failed to revert reservation %s status: %v\n", res.ID, revertErr)
			}
			continue
		}

		log.Printf("Published reservation %s to queue\n", res.ID)
	}

	return nil
}
