package domain

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type ReservationStatus string

const (
	StatusPending    ReservationStatus = "PENDING"
	StatusProcessing ReservationStatus = "PROCESSING"
	StatusDone       ReservationStatus = "DONE"
	StatusFailed     ReservationStatus = "FAILED"
)

type (
	ReservationRepository interface {
		Create(ctx context.Context, reservation Reservation) error                              // 予約を保存する
		GetPendingReservations(ctx context.Context, executeAt time.Time) ([]Reservation, error) // 「待ち状態」で「実行時間を過ぎている」予約を探してくる
		UpdateStatus(ctx context.Context, id string, status ReservationStatus) error            // 予約の状態を更新する
		FindByUserID(ctx context.Context, userID string) ([]Reservation, error)                 // ユーザーIDに紐づく予約一覧を取得する
	}

	Reservation struct {
		ID             string
		UserID         string
		PointAmount    int
		ExecuteAt      time.Time
		Status         ReservationStatus
		IdempotencyKey string
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
)

var (
	ErrInvalidPointAmount       = errors.New("point amount must be greater than 0")
	ErrInvalidExecuteAt         = errors.New("execute_at must be in the future")
	ErrReservationNotFound      = errors.New("reservation not found")
	ErrReservationAlreadyExists = errors.New("reservation with this idempotency key already exists")
	ErrReservationUpdateFailed  = errors.New("failed to update reservation status")
	ErrReservationCreateFailed  = errors.New("failed to create reservation")
	ErrReservationGetFailed     = errors.New("failed to get pending reservations")
)

// 新しい「予約」を作るための関数。
func NewReservation(userID string, pointAmount int, executeAt time.Time) (Reservation, error) {
	if !isCorrectFormatUserID(userID) {
		return Reservation{}, ErrInvalidFormatUserID
	}
	if pointAmount <= 0 {
		return Reservation{}, ErrInvalidPointAmount
	}

	now := time.Now()
	id := uuid.New().String()
	idempotencyKey := id // 予約IDをそのまま冪等性キーとして使用

	return Reservation{
		ID:             id,
		UserID:         userID,
		PointAmount:    pointAmount,
		ExecuteAt:      executeAt,
		Status:         StatusPending,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// 「重複防止キー」を自分で指定したい場合の関数。
func NewReservationWithIdempotencyKey(userID string, pointAmount int, executeAt time.Time, idempotencyKey string) (Reservation, error) {
	if !isCorrectFormatUserID(userID) {
		return Reservation{}, ErrInvalidFormatUserID
	}
	if pointAmount <= 0 {
		return Reservation{}, ErrInvalidPointAmount
	}
	if !isValidIdempotencyKey(idempotencyKey) {
		return Reservation{}, errors.New("idempotency key is invalid format")
	}

	now := time.Now()
	id := uuid.New().String()

	return Reservation{
		ID:             id,
		UserID:         userID,
		PointAmount:    pointAmount,
		ExecuteAt:      executeAt,
		Status:         StatusPending,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// 重複防止キーが変な文字を使ってないか、長すぎないかチェックする。
func isValidIdempotencyKey(key string) bool {
	if len(key) == 0 || len(key) > 100 {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return re.MatchString(key)
}
