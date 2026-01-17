package domain

import (
	"testing"
	"time"
)

func TestNewReservation_Success(t *testing.T) {
	t.Parallel()

	userID := "testUser1"
	pointAmount := 100
	executeAt := time.Now().Add(1 * time.Hour)

	reservation, err := NewReservation(userID, pointAmount, executeAt)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if reservation.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, reservation.UserID)
	}
	if reservation.PointAmount != pointAmount {
		t.Errorf("expected pointAmount %d, got %d", pointAmount, reservation.PointAmount)
	}
	if reservation.Status != StatusPending {
		t.Errorf("expected status %s, got %s", StatusPending, reservation.Status)
	}
	if reservation.ID == "" {
		t.Error("expected non-empty ID")
	}
	if reservation.IdempotencyKey == "" {
		t.Error("expected non-empty IdempotencyKey")
	}
}

func TestNewReservation_InvalidUserID(t *testing.T) {
	t.Parallel()

	_, err := NewReservation("invalid-user!", 100, time.Now().Add(1*time.Hour))
	if err == nil {
		t.Error("expected error for invalid user ID")
	}
	if err != ErrInvalidFormatUserID {
		t.Errorf("expected ErrInvalidFormatUserID, got %v", err)
	}
}

func TestNewReservation_InvalidPointAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		pointAmount int
	}{
		{"zero", 0},
		{"negative", -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewReservation("userA", tt.pointAmount, time.Now().Add(1*time.Hour))
			if err == nil {
				t.Error("expected error for invalid point amount")
			}
			if err != ErrInvalidPointAmount {
				t.Errorf("expected ErrInvalidPointAmount, got %v", err)
			}
		})
	}
}

func TestNewReservationWithIdempotencyKey_Success(t *testing.T) {
	t.Parallel()

	userID := "testUser1"
	pointAmount := 100
	executeAt := time.Now().Add(1 * time.Hour)
	idempotencyKey := "custom-key-123"

	reservation, err := NewReservationWithIdempotencyKey(userID, pointAmount, executeAt, idempotencyKey)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if reservation.IdempotencyKey != idempotencyKey {
		t.Errorf("expected idempotencyKey %s, got %s", idempotencyKey, reservation.IdempotencyKey)
	}
}

func TestNewReservationWithIdempotencyKey_InvalidKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  string
	}{
		{"empty", ""},
		{"special chars", "key@with#special"},
		{"too long", string(make([]byte, 101))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewReservationWithIdempotencyKey("userA", 100, time.Now().Add(1*time.Hour), tt.key)
			if err == nil {
				t.Error("expected error for invalid idempotency key")
			}
		})
	}
}

func TestIsValidIdempotencyKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		key      string
		expected bool
	}{
		{"validkey", true},
		{"valid-key", true},
		{"valid_key", true},
		{"ValidKey123", true},
		{"", false},
		{"key@invalid", false},
		{"key with space", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := isValidIdempotencyKey(tt.key)
			if result != tt.expected {
				t.Errorf("expected %v for key %q, got %v", tt.expected, tt.key, result)
			}
		})
	}
}
