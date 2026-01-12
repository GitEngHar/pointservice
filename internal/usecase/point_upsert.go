package usecase

import (
	"context"
	"errors"
	"pointservice/internal/domain"
	"pointservice/internal/usecase/tally"
	"time"
)

const InitialPointValue = 0

type (
	PointUpsertUseCase interface {
		Execute(context.Context, *PointUpsertInput) error
	}

	PointUpsertInput struct {
		UserID   string `json:"user_id"`
		PointNum int    `json:"point_num"`
	}

	pointAddInterceptor struct {
		repo     domain.PointRepository
		producer tally.Producer
	}
)

func NewPointUpsertInterceptor(
	repo domain.PointRepository,
	producer tally.Producer,
) PointUpsertUseCase {
	return pointAddInterceptor{
		repo:     repo,
		producer: producer,
	}
}

func (p pointAddInterceptor) Execute(ctx context.Context, input *PointUpsertInput) error {
	currentUserPoint, err := p.repo.GetPointByUserID(ctx, input.UserID)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}
	now := time.Now()
	if errors.Is(err, domain.ErrUserNotFound) {
		if currentUserPoint, err = domain.NewPoint(input.UserID, InitialPointValue, now, now); err != nil {
			return err
		}
	}

	updatePoint, err := domain.NewPoint(
		currentUserPoint.UserID,
		currentUserPoint.PointNum+input.PointNum,
		currentUserPoint.CreatedAt,
		now)
	if err != nil {
		return err
	}
	event, err := domain.NewPoint(input.UserID, input.PointNum, now, now)
	if err = p.producer.PublishPoint(ctx, event); err != nil {
		return err
	}
	return p.repo.UpdatePointOrCreateByUserID(ctx, updatePoint)
}
