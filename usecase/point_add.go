package usecase

import (
	"context"
	"pointservice/domain"
)

type (
	PointAddUseCase interface {
		Execute(context.Context, PointAddInput) error
	}

	PointAddInput struct {
		UserID   string `json:"user_id"`
		PointNum int    `json:"point_num"`
	}

	pointAddInterceptor struct {
		repo domain.PointRepository
	}
)

// TODO:
func NewPointAddInterceptor(
	repo domain.PointRepository,
) PointAddUseCase {
	return pointAddInterceptor{
		repo: repo,
	}
}

// TODO: 実行内容を加える UPSERT
func (p pointAddInterceptor) Execute(ctx context.Context, input PointAddInput) error {
	currentUserPoint, err := p.repo.GetPointByUserID(ctx, input.UserID)
	if err != nil {
		return err
	}
	addedPoints, err := domain.NewPoint(currentUserPoint.UserID, currentUserPoint.PointNum+input.PointNum)
	if err != nil {
		return err
	}
	return p.repo.UpdatePointByUserID(ctx, addedPoints) //TODO: そのまま返さない方がいいかもしてない
}
