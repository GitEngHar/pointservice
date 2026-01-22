package usecase

import (
	"context"
	"fmt"
	"log"
	"pointservice/internal/domain"
	"pointservice/internal/infra/aync/dto"
	"pointservice/internal/infra/aync/rabbitmq"
	"pointservice/internal/infra/repository"
	"time"
)

type PointReservationConfirmUseCase interface {
	Execute(ctx context.Context, producer *rabbitmq.RabbitProducer) error
}

type PointReservationConfirmUseCaseImpl struct {
	pointRepo       repository.PointRepository
	reservationRepo repository.ReservationRepository
}

func NewPointReservationConfirmUseCase(
	pointRepo repository.PointRepository,
	reservationRepo repository.ReservationRepository,
) PointReservationConfirmUseCase {
	return &PointReservationConfirmUseCaseImpl{
		pointRepo:       pointRepo,
		reservationRepo: reservationRepo,
	}
}

func (p PointReservationConfirmUseCaseImpl) Execute(ctx context.Context, producer *rabbitmq.RabbitProducer) error {
	now := time.Now()
	reservations, err := p.reservationRepo.GetPendingReservations(ctx, now) // 実行待ちの予約を取得する。
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
		if err := p.reservationRepo.UpdateStatus(ctx, res.ID, domain.StatusProcessing); err != nil {
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
			if revertErr := p.reservationRepo.UpdateStatus(ctx, res.ID, domain.StatusPending); revertErr != nil {
				log.Printf("Failed to revert reservation %s status: %v\n", res.ID, revertErr)
			}
			continue
		}

		log.Printf("Published reservation %s to queue\n", res.ID)
	}

	return nil
}
