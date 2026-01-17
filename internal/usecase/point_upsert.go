package usecase

import (
	"context"
	"errors"
	"pointservice/internal/domain"
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
		repo domain.PointRepository
	}
)

func NewPointUpsertInterceptor(
	repo domain.PointRepository,
) PointUpsertUseCase {
	return pointAddInterceptor{
		repo: repo,
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
	return p.repo.UpdatePointOrCreateByUserID(ctx, updatePoint)
}
