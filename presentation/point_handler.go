package presentation

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"pointservice/adapter/repository"
	"pointservice/domain"
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
	returnMassage := "Success"
	return c.String(http.StatusOK, returnMassage)
}

func handleErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrPointBelowZero):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidFormatUserID):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrUserNotFound):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrSelectUserID):
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	case errors.Is(err, domain.ErrUpdatePoint):
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	case errors.Is(err, domain.ErrCreateOrUpdatePoint):
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	returnMassage := "Success"
	return c.String(http.StatusOK, returnMassage)
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
	returnMassage := fmt.Sprintf("{userID:%s, point:%d}", pointInfo.UserID, pointInfo.PointNum)
	return c.String(http.StatusOK, returnMassage)
}
