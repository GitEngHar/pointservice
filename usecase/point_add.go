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

func NewPointAddInterceptor(
	repo domain.PointRepository,
) PointAddUseCase {
	return pointAddInterceptor{
		repo: repo,
	}
}

// TODO: 実行内容を加える UPSERT
func (p pointAddInterceptor) Execute(ctx context.Context, input PointAddInput) error {
	return nil
}
