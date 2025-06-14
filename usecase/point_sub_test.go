package usecase

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"pointservice/domain"
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
		name              string
		args              args
		expectedErr       bool
		expectedErrorType error
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
			expectedErr:       false,
			expectedErrorType: nil,
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
			expectedErr:       true,
			expectedErrorType: domain.ErrUserNotFound,
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
			expectedErr:       true,
			expectedErrorType: domain.ErrPointBelowZero,
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := uc.Execute(ts.args.Context, ts.args.PointSubInput)
			if ts.expectedErr {
				if diff := cmp.Diff(ts.expectedErrorType, err, cmpopts.EquateErrors()); diff != "" {
					t.Error(diff)
				}
			}
		})
	}

}
