package usecase

import (
	"context"
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

// TODO:
func NewPointSubInterceptor(
	repo domain.PointRepository,
) PointSubUseCase {
	return pointSubInterceptor{
		repo: repo,
	}
}

// TODO: 実行内容を加える UPSERT
func (p pointSubInterceptor) Execute(ctx context.Context, input *PointSubInput) error {
	currentUserPoint, err := p.repo.GetPointByUserID(ctx, input.UserID)
	if err != nil {
		return err
	}
	addedPoints, err := domain.NewPoint(currentUserPoint.UserID, currentUserPoint.PointNum-input.PointNum)
	if err != nil {
		return err
	}
	return p.repo.UpdatePointByUserID(ctx, addedPoints) //TODO: そのまま返さない方がいいかもしてない
}
