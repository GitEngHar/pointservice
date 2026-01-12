package usecase

import (
	"context"
	"pointservice/internal/domain"
	"time"
)

type (
	PointSubUseCase interface {
		Execute(context.Context, *PointSubInput) error
	}

	PointSubInput struct {
		UserID   string `json:"user_id"`
		PointNum int    `json:"point_num"`
	}

	pointSubInterceptor struct {
		repo domain.PointRepository
	}
)

func NewPointSubInterceptor(
	repo domain.PointRepository,
) PointSubUseCase {
	return pointSubInterceptor{
		repo: repo,
	}
}

func (p pointSubInterceptor) Execute(ctx context.Context, input *PointSubInput) error {
	currentUserPoint, err := p.repo.GetPointByUserID(ctx, input.UserID)
	now := time.Now()
	if err != nil {
		return err
	}
	updatePoint, err := domain.NewPoint(
		currentUserPoint.UserID,
		currentUserPoint.PointNum-input.PointNum,
		currentUserPoint.CreatedAt,
		now)
	if err != nil {
		return err
	}
	return p.repo.UpdatePointByUserID(ctx, updatePoint)
}
