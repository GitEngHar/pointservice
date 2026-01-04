package tally

import (
	"context"
)

type PointTally struct {
	consumer Consumer
}

func NewPointTally(consumer Consumer) *PointTally {
	return &PointTally{
		consumer: consumer,
	}
}

func (p *PointTally) Execute(ctx context.Context) error {
	return p.consumer.GetSumPoint(ctx)
}
