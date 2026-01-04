package tally

import "time"

type SumPointSpecifyDate struct {
	Point int       `json:"point"`
	Date  time.Time `json:"date"`
}

type SumPointPerUser struct {
	userID string
	point  int
	date   time.Time
}
