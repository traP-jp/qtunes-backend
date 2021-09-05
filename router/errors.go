package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo/v4"
)

func generateEchoError(err error) error {
	if errors.Is(err, model.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Not Found")
	} else if errors.Is(err, model.ErrNoChange) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "No Change")
	} else {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
}

func errSessionNotFound(err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Failed in Getting Session:%w", err).Error())
}

func errBind(err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Failed to bind request: %w", err).Error())
}
