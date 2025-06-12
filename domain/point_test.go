package domain

import (
	"github.com/google/go-cmp/cmp"
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
		name     string
		args     args
		expected Point
		errMsg   string
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
			errMsg: "",
		},
		{
			name: "Failed create point struct",
			args: args{
				UserID:    "!abc123d",
				PointNum:  1000,
				CreatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
			},
			expected: Point{},
			errMsg:   "userID is not correct format: !abc123d",
		},
		{
			name: "Failed create point struct",
			args: args{
				UserID:    "abc123d",
				PointNum:  -1,
				CreatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, time.June, 10, 15, 30, 0, 0, time.UTC),
			},
			expected: Point{},
			errMsg:   "points must be greater than 0",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			result, err := NewPoint(ts.args.UserID, ts.args.PointNum, ts.args.CreatedAt, ts.args.UpdatedAt)
			if diff := cmp.Diff(ts.expected, result); diff != "" {
				t.Error(diff)
			}
			var errMsg string
			if err != nil {
				errMsg = err.Error()
				if diff := cmp.Diff(ts.errMsg, errMsg); diff != "" {
					t.Error(diff)
				}
			}
		})
	}

}
