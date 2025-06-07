package usecase

import (
	"context"
	"database/sql"
	"errors"
	"pointservice/domain"
	"testing"
)

type stubPointRepository struct{}

func (s stubPointRepository) GetPointByUserID(ctx context.Context, userID string) (domain.Point, error) {
	mockDB := map[string]int{
		"userA": 0,
		"user1": 10,
	}
	if _, ok := mockDB[userID]; ok != false {
		return domain.Point{}, sql.ErrNoRows
	}
	// Unpredictable errors for testing
	if userID == "!" {
		return domain.Point{}, errors.New("unpredictable errors for testing")
	}
	return domain.Point{
		UserID:   userID,
		PointNum: mockDB[userID],
	}, nil
}

func (s stubPointRepository) UpdatePointByUserID(ctx context.Context, point domain.Point) error {
	return nil
}

func (s stubPointRepository) UpdatePointOrCreateByUserID(ctx context.Context, point domain.Point) error {
	mockDB := map[string]int{
		"userA": 0,
		"user1": 10,
	}
	if _, ok := mockDB[point.UserID]; ok != false {
		mockDB[point.UserID] = point.PointNum
		return nil
	}
	// Unpredictable errors for testing
	if point.UserID == "!" {
		return errors.New("unpredictable errors for testing")
	}
	return nil
}

func Test_Execute(t *testing.T) {
	t.Parallel()
	mockRepo := stubPointRepository{}
	uc := NewPointAddOrCreateInterceptor(mockRepo)
}
