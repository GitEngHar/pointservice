package usecase

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"testing"
)

// Test_Execute name is deprecated add or create uc
func Test_Sub_Execute(t *testing.T) {
	type args struct {
		context.Context
		*PointSubInput
	}

	t.Parallel()

	mockRepo := StubPointRepository{}
	uc := NewPointSubInterceptor(mockRepo)
	tests := []struct {
		name        string
		args        args
		expectedErr bool
		errMsg      string
	}{
		{
			name: "Successful point sub",
			args: args{
				context.Background(),
				&PointSubInput{
					"userB",
					9,
				},
			},
			expectedErr: false,
			errMsg:      "",
		},
		{
			name: "Failed select target user:",
			args: args{
				context.Background(),
				&PointSubInput{
					"userC",
					9,
				},
			},
			expectedErr: true,
			errMsg:      "failed select target user: sql: no rows in result set",
		},
		{
			name: "Failed generate new point",
			args: args{
				context.Background(),
				&PointSubInput{
					"userA",
					1,
				},
			},
			expectedErr: true,
			errMsg:      "point update failed: points must be greater than 0",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := uc.Execute(ts.args.Context, ts.args.PointSubInput)
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
