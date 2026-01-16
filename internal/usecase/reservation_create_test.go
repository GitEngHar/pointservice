package usecase

import (
	"context"
	"pointservice/internal/domain"
	"testing"
	"time"
)

type StubReservationRepository struct {
	createCalled bool
	lastCreated  domain.Reservation
	createErr    error
}

func (s *StubReservationRepository) Create(ctx context.Context, reservation domain.Reservation) error {
	s.createCalled = true
	s.lastCreated = reservation
	return s.createErr
}

func (s *StubReservationRepository) GetPendingReservations(ctx context.Context, executeAt time.Time) ([]domain.Reservation, error) {
	return nil, nil
}

func (s *StubReservationRepository) UpdateStatus(ctx context.Context, id string, status domain.ReservationStatus) error {
	return nil
}

func TestReservationCreate_Success(t *testing.T) {
	t.Parallel()

	repo := &StubReservationRepository{}
	uc := NewReservationCreateInterceptor(repo)

	input := &ReservationCreateInput{
		UserID:      "testUser1",
		PointAmount: 100,
		ExecuteAt:   time.Now().Add(1 * time.Hour),
	}

	output, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if output == nil {
		t.Fatal("expected output, got nil")
	}
	if output.ReservationID == "" {
		t.Error("expected non-empty reservation ID")
	}
	if output.Status != string(domain.StatusPending) {
		t.Errorf("expected status %s, got %s", domain.StatusPending, output.Status)
	}
	if !repo.createCalled {
		t.Error("expected repository Create to be called")
	}
}

func TestReservationCreate_WithIdempotencyKey(t *testing.T) {
	t.Parallel()

	repo := &StubReservationRepository{}
	uc := NewReservationCreateInterceptor(repo)

	input := &ReservationCreateInput{
		UserID:         "testUser1",
		PointAmount:    100,
		ExecuteAt:      time.Now().Add(1 * time.Hour),
		IdempotencyKey: "custom-key-123",
	}

	output, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if repo.lastCreated.IdempotencyKey != input.IdempotencyKey {
		t.Errorf("expected idempotency key %s, got %s", input.IdempotencyKey, repo.lastCreated.IdempotencyKey)
	}
	if output.ReservationID == "" {
		t.Error("expected non-empty reservation ID")
	}
}

func TestReservationCreate_InvalidUserID(t *testing.T) {
	t.Parallel()

	repo := &StubReservationRepository{}
	uc := NewReservationCreateInterceptor(repo)

	input := &ReservationCreateInput{
		UserID:      "invalid-user!",
		PointAmount: 100,
		ExecuteAt:   time.Now().Add(1 * time.Hour),
	}

	_, err := uc.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error for invalid user ID")
	}
	if repo.createCalled {
		t.Error("repository Create should not be called on validation error")
	}
}

func TestReservationCreate_InvalidPointAmount(t *testing.T) {
	t.Parallel()

	repo := &StubReservationRepository{}
	uc := NewReservationCreateInterceptor(repo)

	input := &ReservationCreateInput{
		UserID:      "testUser1",
		PointAmount: 0,
		ExecuteAt:   time.Now().Add(1 * time.Hour),
	}

	_, err := uc.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error for zero point amount")
	}
	if repo.createCalled {
		t.Error("repository Create should not be called on validation error")
	}
}

func TestReservationCreate_CreateError(t *testing.T) {
	t.Parallel()

	repo := &StubReservationRepository{
		createErr: domain.ErrReservationCreateFailed,
	}
	uc := NewReservationCreateInterceptor(repo)

	input := &ReservationCreateInput{
		UserID:      "testUser1",
		PointAmount: 100,
		ExecuteAt:   time.Now().Add(1 * time.Hour),
	}

	_, err := uc.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error from repository")
	}
}
