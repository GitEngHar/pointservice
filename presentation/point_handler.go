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
	fmt.Println("Recieved PointAdd")
	ctx := c.Request().Context()
	pointDTO := new(usecase.PointAddInput)
	if err := c.Bind(pointDTO); err != nil {
		return err
	}
	fmt.Println("Bind json request")
	uc := usecase.NewPointAddInterceptor(p.repo)
	if err := uc.Execute(ctx, pointDTO); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("point added")
	returnMassage := "Success"
	return c.String(http.StatusOK, returnMassage)
}

func (p *PointHandler) PointSub(c echo.Context) error {
	returnMassage := "pointSub"
	return c.String(http.StatusOK, returnMassage)
}
