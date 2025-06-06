package usecase

import (
	"context"
	"database/sql"
	"errors"
	"pointservice/domain"
)

type (
	PointAddUseCase interface {
		Execute(context.Context, *PointAddInput) error
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

func (p pointAddInterceptor) Execute(ctx context.Context, input *PointAddInput) error {
	currentUserPoint, err := p.repo.GetPointByUserID(ctx, input.UserID)
	if err != nil {
		// Sql.ErrNoRows are tolerated as new users are created.
		if errors.Is(err, sql.ErrNoRows) {
			// The parameters of the user to be created refer to the input parameters
			currentUserPoint.UserID = input.UserID
			currentUserPoint.PointNum = 0
		} else {
			return err
		}
	}
	addedPoints, err := domain.NewPoint(currentUserPoint.UserID, currentUserPoint.PointNum+input.PointNum)
	if err != nil {
		return err
	}
	return p.repo.UpdatePointOrCreateByUserID(ctx, addedPoints)
}
