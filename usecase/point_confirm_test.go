package usecase

import (
	"context"
	"github.com/google/go-cmp/cmp"
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
		name        string
		args        args
		expected    domain.Point
		expectedErr bool
		errMsg      string
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
			expectedErr: false,
			errMsg:      "",
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
			expectedErr: false,
			errMsg:      "",
		},
		{
			name: "Failed select target user:",
			args: args{
				context.Background(),
				&PointConfirmInput{
					"userC",
				},
			},
			expected:    domain.Point{},
			expectedErr: true,
			errMsg:      "failed select target user: sql: no rows in result set",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			pointInfo, err := uc.Execute(ts.args.Context, ts.args.PointConfirmInput)
			if diff := cmp.Diff(ts.expected, pointInfo); diff != "" {
				t.Error(diff)
			}
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
