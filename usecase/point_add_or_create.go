package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pointservice/domain"
)

type (
	PointAddOrCreateUseCase interface {
		Execute(context.Context, *PointAddOrCreateInput) error
	}

	PointAddOrCreateInput struct {
		UserID   string `json:"user_id"`
		PointNum int    `json:"point_num"`
	}

	pointAddInterceptor struct {
		repo domain.PointRepository
	}
)

func NewPointAddOrCreateInterceptor(
	repo domain.PointRepository,
) PointAddOrCreateUseCase {
	return pointAddInterceptor{
		repo: repo,
	}
}

func (p pointAddInterceptor) Execute(ctx context.Context, input *PointAddOrCreateInput) error {
	currentUserPoint, err := p.repo.GetPointByUserID(ctx, input.UserID)
	if err != nil {
		// Sql.ErrNoRows are tolerated as new users are created.
		if errors.Is(err, sql.ErrNoRows) {
			// The parameters of the user to be created refer to the input parameters
			currentUserPoint.UserID = input.UserID
			currentUserPoint.PointNum = 0
		} else {
			return fmt.Errorf("failed select target user: %w", err)
		}
	}
	addedPoints, err := domain.NewPoint(currentUserPoint.UserID, currentUserPoint.PointNum+input.PointNum)
	if err != nil {
		return fmt.Errorf("new point create filed: %w", err)
	}
	return p.repo.UpdatePointOrCreateByUserID(ctx, addedPoints)
}
