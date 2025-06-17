package presentation

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"pointservice/adapter/repository"
	"pointservice/domain"
	"pointservice/presentation/api"
	"pointservice/usecase"
)

type PointHandler struct {
	db   *sql.DB
	repo repository.PointRepository
}

func NewPointHandler(db *sql.DB, pointRepository repository.PointRepository) *PointHandler {
	return &PointHandler{
		db:   db,
		repo: pointRepository}
}

func (p *PointHandler) HealthCheck(c echo.Context) error {
	successMessage := api.NewSuccess([]string{"ok"})
	return c.JSON(http.StatusOK, successMessage)
}

func (p *PointHandler) PointAdd(c echo.Context) error {
	ctx := c.Request().Context()
	pointDTO := new(usecase.PointAddOrCreateInput)
	if err := c.Bind(pointDTO); err != nil {
		return handleErr(err)
	}
	uc := usecase.NewPointAddOrCreateInterceptor(p.repo)
	if err := uc.Execute(ctx, pointDTO); err != nil {
		return handleErr(err)
	}
	successMessage := api.NewSuccess([]string{"point updated"})
	return c.JSON(http.StatusOK, successMessage)
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
