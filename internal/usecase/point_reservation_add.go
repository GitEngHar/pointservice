package usecase

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"pointservice/internal/domain"
	"pointservice/internal/infra/aync/dto"
	"pointservice/internal/infra/repository"
)

type PointReservationAddUseCase interface {
	Execute(ctx context.Context, msg amqp.Delivery) error
}
type PointReservationAddUseCaseImpl struct {
	pointRepo       repository.PointRepository
	reservationRepo repository.ReservationRepository
}

func NewPointReservationAddUseCase(pointRepo repository.PointRepository, reservationRepo repository.ReservationRepository) PointReservationAddUseCase {
	return &PointReservationAddUseCaseImpl{
		pointRepo:       pointRepo,
		reservationRepo: reservationRepo,
	}
}

// Execute 予約メッセージを処理する
func (a PointReservationAddUseCaseImpl) Execute(
	ctx context.Context,
	msg amqp.Delivery,
) error {
	var reservation dto.ReservationMessage
	if err := json.Unmarshal(msg.Body, &reservation); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return err
	}

	log.Printf("Processing reservation: %s (user: %s, amount: %d)",
		reservation.ReservationID, reservation.UserID, reservation.PointAmount)

	// Idempotent point grant
	added, err := a.pointRepo.AddPointIdempotent(
		ctx,
		reservation.IdempotencyKey,
		reservation.UserID,
		reservation.PointAmount,
	)
	if err != nil {
		log.Printf("Failed to add point for reservation %s: %v", reservation.ReservationID, err)
		// Update reservation status to FAILED
		if updateErr := a.reservationRepo.UpdateStatus(ctx, reservation.ReservationID, domain.StatusFailed); updateErr != nil {
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
	if err := a.reservationRepo.UpdateStatus(ctx, reservation.ReservationID, domain.StatusDone); err != nil {
		log.Printf("Failed to update reservation status to DONE: %v", err)
		return err
	}

	log.Printf("Reservation %s completed", reservation.ReservationID)
	return nil
}
