package domain

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func Test_NewPoint(t *testing.T) {
	t.Parallel()
	type args struct {
		UserID   string
		PointNum int
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
				UserID:   "abc123d",
				PointNum: 1000,
			},
			expected: Point{
				UserID:   "abc123d",
				PointNum: 1000,
			},
			errMsg: "",
		},
		{
			name: "Failed create point struct",
			args: args{
				UserID:   "!abc123d",
				PointNum: 1000,
			},
			expected: Point{},
			errMsg:   "userID is not correct format: !abc123d",
		},
		{
			name: "Failed create point struct",
			args: args{
				UserID:   "abc123d",
				PointNum: -1,
			},
			expected: Point{},
			errMsg:   "points must be greater than 0",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			result, err := NewPoint(ts.args.UserID, ts.args.PointNum)
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
