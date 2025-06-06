package usecase

import (
	"context"
	"fmt"
	"pointservice/domain"
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
	if err != nil {
		return fmt.Errorf("failed select target user: %w", err)
	}
	subtractedPoints, err := domain.NewPoint(currentUserPoint.UserID, currentUserPoint.PointNum-input.PointNum)
	if err != nil {
		return fmt.Errorf("new point create filed: %w", err)
	}
	return p.repo.UpdatePointByUserID(ctx, subtractedPoints)
}
