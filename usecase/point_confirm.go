package usecase

import (
	"context"
	"pointservice/domain"
)

type (
	PointConfirmUseCase interface {
		Execute(context.Context, *PointConfirmInput) (domain.Point, error)
	}

	PointConfirmInput struct {
		UserID string `json:"user_id"`
	}

	pointConfirmInterceptor struct {
		repo domain.PointRepository
	}
)

func NewPointConfirmInterceptor(
	repo domain.PointRepository,
) PointConfirmUseCase {
	return pointConfirmInterceptor{
		repo: repo,
	}
}

func (p pointConfirmInterceptor) Execute(ctx context.Context, input *PointConfirmInput) (domain.Point, error) {
	currentUserPoint, err := p.repo.GetPointByUserID(ctx, input.UserID)
	if err != nil {
		return domain.Point{}, err
	}
	pointInfo, err := domain.NewPoint(currentUserPoint.UserID, currentUserPoint.PointNum, currentUserPoint.CreatedAt, currentUserPoint.UpdatedAt)
	if err != nil {
		return domain.Point{}, err
	}
	return pointInfo, nil
}
