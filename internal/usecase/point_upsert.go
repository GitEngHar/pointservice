package usecase

import (
	"context"
	"errors"
	"pointservice/internal/domain"
	"pointservice/internal/usecase/tally"
	"time"
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
		repo     domain.PointRepository
		producer tally.Producer
	}
)

func NewPointAddOrCreateInterceptor(
	repo domain.PointRepository,
	producer tally.Producer,
) PointAddOrCreateUseCase {
	return pointAddInterceptor{
		repo:     repo,
		producer: producer,
	}
}

func (p pointAddInterceptor) Execute(ctx context.Context, input *PointAddOrCreateInput) error {
	currentUserPoint, err := p.repo.GetPointByUserID(ctx, input.UserID)
	if err != nil {
		// Sql.ErrNoRows are tolerated as new users are created.
		if errors.Is(err, domain.ErrUserNotFound) {
			// The parameters of the user to be created refer to the input parameters
			if currentUserPoint, err = domain.NewPoint(input.UserID, 0, time.Now(), time.Now()); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	addedPoints, err := domain.NewPoint(currentUserPoint.UserID, currentUserPoint.PointNum+input.PointNum, currentUserPoint.CreatedAt, time.Now())
	if err != nil {
		return err
	}
	addPoint, err := domain.NewPoint(input.UserID, input.PointNum, time.Now(), time.Now())
	if err = p.producer.PublishPoint(ctx, addPoint); err != nil {
		return err
	}
	return p.repo.UpdatePointOrCreateByUserID(ctx, addedPoints)
}
