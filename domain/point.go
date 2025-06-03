package domain

import (
	"errors"
	"regexp"
)

type Point struct {
	userID   string
	pointNum int
}

func isCorrectFormatUserID(target string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(target)
}

func NewPoint(userID string, pointNum int) (*Point, error) {
	if !isCorrectFormatUserID(userID) {
		return nil, errors.New("userID is not correct format :" + userID)
	}
	if 0 > pointNum {
		return nil, errors.New("points must be greater than 0")
	}
	return &Point{
		userID:   userID,
		pointNum: pointNum,
	}, nil
}
