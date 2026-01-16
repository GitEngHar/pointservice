package usecase

import (
	"context"
	"pointservice/internal/domain"
	"time"
)

type (

	// 「Execute」というボタンを押すと、予約作成の処理が走る約束になっている。
	ReservationCreateUseCase interface {
		Execute(context.Context, *ReservationCreateInput) (*ReservationCreateOutput, error)
	}

	// 予約を作るときに必要な「材料」を入れるための箱。
	ReservationCreateInput struct {
		UserID         string    `json:"user_id"`
		PointAmount    int       `json:"point_amount"`
		ExecuteAt      time.Time `json:"execute_at"`
		IdempotencyKey string    `json:"idempotency_key,omitempty"`
	}

	// 予約が終わった後に返す「結果」を入れるための箱。
	ReservationCreateOutput struct {
		ReservationID string    `json:"reservation_id"`
		Status        string    `json:"status"`
		ExecuteAt     time.Time `json:"execute_at"`
	}

	// 実際に予約作成の仕事をする「担当者（構造体）」。
	reservationCreateInterceptor struct {
		repo domain.ReservationRepository
	}
)

func NewReservationCreateInterceptor(repo domain.ReservationRepository) ReservationCreateUseCase {
	return reservationCreateInterceptor{
		repo: repo,
	}
}

func (r reservationCreateInterceptor) Execute(ctx context.Context, input *ReservationCreateInput) (*ReservationCreateOutput, error) {
	var reservation domain.Reservation
	var err error

	if input.IdempotencyKey != "" {
		reservation, err = domain.NewReservationWithIdempotencyKey(
			input.UserID,
			input.PointAmount,
			input.ExecuteAt,
			input.IdempotencyKey,
		)
	} else {

		// 新しい予約データを作る。
		reservation, err = domain.NewReservation(
			input.UserID,
			input.PointAmount,
			input.ExecuteAt,
		)
	}
	if err != nil {
		return nil, err
	}

	// データベース係（r.repo）にお願いして、作った予約データを保存（Create）してもらう。
	if err := r.repo.Create(ctx, reservation); err != nil {
		return nil, err
	}

	return &ReservationCreateOutput{
		ReservationID: reservation.ID,
		Status:        string(reservation.Status),
		ExecuteAt:     reservation.ExecuteAt,
	}, nil
}
