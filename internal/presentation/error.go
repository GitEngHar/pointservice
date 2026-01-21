package presentation

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"pointservice/internal/domain"
	"pointservice/internal/presentation/api"
)

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
