package presentation

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"pointservice/internal/domain"
	"pointservice/internal/infra/repository"
	"pointservice/internal/presentation/api"
	"pointservice/internal/usecase"

	"github.com/labstack/echo/v4"
)

type PointHandler struct {
	db              *sql.DB
	repo            repository.PointRepository
	reservationRepo repository.ReservationRepository
}

func NewPointHandler(db *sql.DB, pointRepository repository.PointRepository, reservationRepository repository.ReservationRepository) *PointHandler {
	return &PointHandler{
		db:              db,
		repo:            pointRepository,
		reservationRepo: reservationRepository,
	}
}

func (p *PointHandler) HealthCheck(c echo.Context) error {
	successMessage := api.NewSuccess([]string{"ok"})
	return c.JSON(http.StatusOK, successMessage)
}

func (p *PointHandler) PointAdd(c echo.Context) error {
	ctx := c.Request().Context()
	pointDTO := new(usecase.PointUpsertInput)
	if err := c.Bind(pointDTO); err != nil {
		return handleErr(err)
	}
	uc := usecase.NewPointUpsertInterceptor(p.repo)
	if err := uc.Execute(ctx, pointDTO); err != nil {
		return handleErr(err)
	}
	successMessage := api.NewSuccess([]string{"point updated"})
	return c.JSON(http.StatusOK, successMessage)
}

func (p *PointHandler) PointSub(c echo.Context) error {
	ctx := c.Request().Context()
	pointDTO := new(usecase.PointSubInput)
	if err := c.Bind(pointDTO); err != nil {
		return handleErr(err)
	}
	uc := usecase.NewPointSubInterceptor(p.repo)
	if err := uc.Execute(ctx, pointDTO); err != nil {
		return handleErr(err)
	}
	successMessage := api.NewSuccess([]string{"point subtracted"})
	return c.JSON(http.StatusOK, successMessage)
}

func (p *PointHandler) PointConfirm(c echo.Context) error {
	ctx := c.Request().Context()
	pointDTO := new(usecase.PointConfirmInput)
	if err := c.Bind(pointDTO); err != nil {
		return handleErr(err)
	}
	uc := usecase.NewPointConfirmInterceptor(p.repo)
	pointInfo, err := uc.Execute(ctx, pointDTO)
	if err != nil {
		return handleErr(err)
	}
	successMessage := api.NewSuccess([]string{
		fmt.Sprintf("userID: %s", pointInfo.UserID),
		fmt.Sprintf("pointNum: %d", pointInfo.PointNum),
	})
	return c.JSON(http.StatusOK, successMessage)
}

// ポイント予約（仮押さえ）を受け付ける窓口となる関数。
func (p *PointHandler) PointReserve(c echo.Context) error {
	ctx := c.Request().Context()
	reservationDTO := new(usecase.ReservationCreateInput) // 予約情報を入れるための「空っぽの箱（構造体）」を作っている。
	if err := c.Bind(reservationDTO); err != nil {
		return handleErr(err)
	}
	uc := usecase.NewReservationCreateInterceptor(p.reservationRepo) // 実際に「予約を作成する仕事」を担当する人（Interceptor）を呼び出して準備している。
	result, err := uc.Execute(ctx, reservationDTO)
	if err != nil {
		return handleErr(err)
	}
	return c.JSON(http.StatusCreated, result)
}

func handleErr(err error) error {
	errMessages := api.NewError(err)
	switch {
	case errors.Is(err, domain.ErrPointBelowZero):
		return echo.NewHTTPError(http.StatusBadRequest, errMessages)
	case errors.Is(err, domain.ErrInvalidFormatUserID):
		return echo.NewHTTPError(http.StatusBadRequest, errMessages)
	case errors.Is(err, domain.ErrUserNotFound):
		return echo.NewHTTPError(http.StatusBadRequest, errMessages)
	case errors.Is(err, domain.ErrSelectUserID):
		return echo.NewHTTPError(http.StatusInternalServerError, errMessages)
	case errors.Is(err, domain.ErrUpdatePoint):
		return echo.NewHTTPError(http.StatusInternalServerError, errMessages)
	case errors.Is(err, domain.ErrCreateOrUpdatePoint):
		return echo.NewHTTPError(http.StatusInternalServerError, errMessages)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, errMessages)
	}
}
