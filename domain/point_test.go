package domain

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	"time"
)

func Test_NewPoint(t *testing.T) {
	t.Parallel()
	type args struct {
		UserID    string
		PointNum  int
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name              string
		args              args
		expected          Point
		expectedErrorType error
	}{
		{
			name: "Successful create point struct",
			args: args{
				UserID:    "abc123d",
				PointNum:  1000,
				CreatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
			},
			expected: Point{
				UserID:    "abc123d",
				PointNum:  1000,
				CreatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
			},
			expectedErrorType: nil,
		},
		{
			name: "Failed create point struct",
			args: args{
				UserID:    "!abc123d",
				PointNum:  1000,
				CreatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
			},
			expected:          Point{},
			expectedErrorType: ErrInvalidFormatUserID,
		},
		{
			name: "Failed create point struct",
			args: args{
				UserID:    "abc123d",
				PointNum:  -1,
				CreatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
			},
			expected:          Point{},
			expectedErrorType: ErrPointBelowZero,
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			result, err := NewPoint(ts.args.UserID, ts.args.PointNum, ts.args.CreatedAt, ts.args.UpdatedAt)
			if diff := cmp.Diff(ts.expected, result); diff != "" {
				t.Error(diff)
			}
			if err != nil {
				if diff := cmp.Diff(ts.expectedErrorType, err, cmpopts.EquateErrors()); diff != "" {
					t.Error(diff)
				}
			}
		})
	}

}
