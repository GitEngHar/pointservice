package domain

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
)

type (
	PointRepository interface {
		GetPointByUserID(ctx context.Context, userID string) (Point, error)
		UpdatePointByUserID(ctx context.Context, point Point) error
		UpdatePointOrCreateByUserID(ctx context.Context, point Point) error
	}

	Point struct {
		UserID    string
		PointNum  int
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

var (
	ErrPointBelowZero = errors.New("points must be greater than 0")
)

func isCorrectFormatUserID(target string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(target)
}

func NewPoint(userID string, pointNum int, createdAt time.Time, updatedAt time.Time) (Point, error) {
	if !isCorrectFormatUserID(userID) {
		return Point{}, fmt.Errorf("userID is not correct format: %s", userID)
	}
	if 0 > pointNum {
		return Point{}, ErrPointBelowZero
	}
	return Point{
		UserID:    userID,
		PointNum:  pointNum,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
