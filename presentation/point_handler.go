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
	fmt.Println("Received PointAdd Request")
	ctx := c.Request().Context()
	pointDTO := new(usecase.PointAddOrCreateInput)
	if err := c.Bind(pointDTO); err != nil {
		return fmt.Errorf("json format bind err: %w", err)
	}
	uc := usecase.NewPointAddOrCreateInterceptor(p.repo)
	if err := uc.Execute(ctx, pointDTO); err != nil {
		return fmt.Errorf("point add or one user create error: %w", err)
	}
	fmt.Println("Success point added")
	returnMassage := "Success"
	return c.String(http.StatusOK, returnMassage)
}

func (p *PointHandler) PointSub(c echo.Context) error {
	fmt.Println("Received PointSub Request")
	ctx := c.Request().Context()
	pointDTO := new(usecase.PointSubInput)
	if err := c.Bind(pointDTO); err != nil {
		return err
	}
	uc := usecase.NewPointSubInterceptor(p.repo)
	if err := uc.Execute(ctx, pointDTO); err != nil {
		fmt.Println(err)
		return fmt.Errorf("point subtraction error: %w", err)
	}
	fmt.Println("Success point subtraction")
	returnMassage := "Success"
	return c.String(http.StatusOK, returnMassage)
}
