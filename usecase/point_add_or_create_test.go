package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/go-cmp/cmp"
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
	type args struct {
		context.Context
		*PointAddOrCreateInput
	}

	t.Parallel()

	mockRepo := stubPointRepository{}
	uc := NewPointAddOrCreateInterceptor(mockRepo)
	tests := []struct {
		name        string
		args        args
		expectedErr bool
		errMsg      string
	}{
		{
			name: "Successful point add",
			args: args{
				context.Background(),
				&PointAddOrCreateInput{
					"userA",
					123,
				},
			},
			expectedErr: false,
			errMsg:      "",
		},
		{
			name: "Successful point create",
			args: args{
				context.Background(),
				&PointAddOrCreateInput{
					"123test",
					123,
				},
			},
			expectedErr: false,
			errMsg:      "",
		},
		{
			name: "Failed select target user",
			args: args{
				context.Background(),
				&PointAddOrCreateInput{
					"!",
					123,
				},
			},
			expectedErr: true,
			errMsg:      "failed select target user: unpredictable errors for testing",
		},
		{
			name: "Failed generate new point",
			args: args{
				context.Background(),
				&PointAddOrCreateInput{
					"123test",
					-1,
				},
			},
			expectedErr: true,
			errMsg:      "new point create filed: points must be greater than 0",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := uc.Execute(ts.args.Context, ts.args.PointAddOrCreateInput)
			var errMsg string
			if ts.expectedErr {
				errMsg = err.Error()
				if diff := cmp.Diff(ts.errMsg, errMsg); diff != "" {
					t.Error(diff)
				}
			}
		})
	}

}
