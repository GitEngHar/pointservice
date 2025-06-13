package presentation

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"pointservice/adapter/repository"
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
		return echo.NewHTTPError(http.StatusBadRequest, "json format bind error: "+err.Error())
	}
	uc := usecase.NewPointAddOrCreateInterceptor(p.repo)
	if err := uc.Execute(ctx, pointDTO); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "point add or one user create error: "+err.Error())
	}
	returnMassage := "Success"
	return c.String(http.StatusOK, returnMassage)
}

func (p *PointHandler) PointSub(c echo.Context) error {
	ctx := c.Request().Context()
	pointDTO := new(usecase.PointSubInput)
	if err := c.Bind(pointDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "json format bind error: "+err.Error())
	}
	uc := usecase.NewPointSubInterceptor(p.repo)
	if err := uc.Execute(ctx, pointDTO); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "point subtraction error: "+err.Error())
	}
	returnMassage := "Success"
	return c.String(http.StatusOK, returnMassage)
}

func (p *PointHandler) PointConfirm(c echo.Context) error {
	ctx := c.Request().Context()
	pointDTO := new(usecase.PointConfirmInput)
	if err := c.Bind(pointDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "json format bind error: "+err.Error())
	}
	uc := usecase.NewPointConfirmInterceptor(p.repo)
	pointInfo, err := uc.Execute(ctx, pointDTO)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "point confirm error: "+err.Error())
	}
	returnMassage := fmt.Sprintf("{userID:%s, point:%d}", pointInfo.UserID, pointInfo.PointNum)
	return c.String(http.StatusOK, returnMassage)
}
