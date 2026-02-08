package usecase

import (
	"context"
	"fmt"
	"pointservice/internal/domain"
	"time"
)

type (
	// Usecase
	ReservationListUsecase struct {
		repo domain.ReservationRepository
	}

	// Input
	ReservationListInput struct {
		UserID string
	}

	// Output
	ReservationListOutput struct {
		ReservedPoints []ReservedPointDTO
	}

	ReservedPointDTO struct {
		Point   int
		Status  string
		AddDate string // RFC3339 format
	}
)

func NewReservationListUsecase(repo domain.ReservationRepository) *ReservationListUsecase {
	return &ReservationListUsecase{
		repo: repo,
	}
}

func (u *ReservationListUsecase) Execute(ctx context.Context, input *ReservationListInput) (*ReservationListOutput, error) {
	if input == nil || input.UserID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	// 1. Repositoryからデータを取得
	reservations, err := u.repo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}

	// 2. DTOに変換
	dtos := make([]ReservedPointDTO, 0, len(reservations))
	for _, r := range reservations {
		dtos = append(dtos, ReservedPointDTO{
			Point:   r.PointAmount,
			Status:  string(r.Status),
			AddDate: r.ExecuteAt.Format(time.RFC3339),
		})
	}

	return &ReservationListOutput{
		ReservedPoints: dtos,
	}, nil
}
