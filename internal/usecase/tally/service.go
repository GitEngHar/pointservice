package tally

import (
	"context"
	"pointservice/internal/domain"
)

type Producer interface {
	PublishPoint(ctx context.Context, point domain.Point) error
}

type Consumer interface {
	GetSumPointSpecifyDate(ctx context.Context) error
}
