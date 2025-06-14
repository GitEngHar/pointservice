package usecase

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"pointservice/domain"
	"testing"
)

// Test_Execute name is deprecated add or create uc
func Test_Confirm_Execute(t *testing.T) {
	type args struct {
		context.Context
		*PointConfirmInput
	}

	t.Parallel()

	mockRepo := StubPointRepository{}
	uc := NewPointConfirmInterceptor(mockRepo)
	tests := []struct {
		name              string
		args              args
		expected          domain.Point
		expectedErr       bool
		expectedErrorType error
	}{
		{
			name: "Successful point confirm",
			args: args{
				context.Background(),
				&PointConfirmInput{
					"userA",
				},
			},
			expected: domain.Point{
				UserID:   "userA",
				PointNum: 0,
			},
			expectedErr:       false,
			expectedErrorType: nil,
		},
		{
			name: "Successful point confirm",
			args: args{
				context.Background(),
				&PointConfirmInput{
					"userB",
				},
			},
			expected: domain.Point{
				UserID:   "userB",
				PointNum: 10,
			},
			expectedErr:       false,
			expectedErrorType: nil,
		},
		{
			name: "Failed select target user:",
			args: args{
				context.Background(),
				&PointConfirmInput{
					"userC",
				},
			},
			expected:          domain.Point{},
			expectedErr:       true,
			expectedErrorType: domain.ErrUserNotFound,
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			pointInfo, err := uc.Execute(ts.args.Context, ts.args.PointConfirmInput)
			if diff := cmp.Diff(ts.expected, pointInfo); diff != "" {
				t.Error(diff)
			}
			if ts.expectedErr {
				if diff := cmp.Diff(ts.expectedErrorType, err, cmpopts.EquateErrors()); diff != "" {
					t.Error(diff)
				}
			}
		})
	}

}
