package domain

import (
	"context"
	"errors"
	"regexp"
)

type (
	PointRepository interface {
		GetPointByUserID(ctx context.Context, userID string) (Point, error)
		UpdatePointByUserID(ctx context.Context, point Point) error
		UpdatePointOrCreateByUserID(ctx context.Context, point Point) error
	}

	Point struct {
		UserID   string
		PointNum int
	}
)

func isCorrectFormatUserID(target string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(target)
}

func NewPoint(userID string, pointNum int) (Point, error) {
	if !isCorrectFormatUserID(userID) {
		return Point{}, errors.New("userID is not correct format :" + userID)
	}
	if 0 > pointNum {
		return Point{}, errors.New("points must be greater than 0")
	}
	return Point{
		UserID:   userID,
		PointNum: pointNum,
	}, nil
}
