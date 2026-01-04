package domain

import (
	"context"
	"errors"
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
		UserID    string    `json:"user_id"`
		PointNum  int       `json:"point_num"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

var (
	ErrPointBelowZero      = errors.New("points must be greater than 0")
	ErrInvalidFormatUserID = errors.New("user id is invalid format")
	ErrUserNotFound        = errors.New("user not found")
	ErrSelectUserID        = errors.New("failed select user")
	ErrUpdatePoint         = errors.New("failed point update")
	ErrCreateOrUpdatePoint = errors.New("failed point create or update")
)

func isCorrectFormatUserID(target string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(target)
}

func NewPoint(userID string, pointNum int, createdAt time.Time, updatedAt time.Time) (Point, error) {
	if !isCorrectFormatUserID(userID) {
		return Point{}, ErrInvalidFormatUserID
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
