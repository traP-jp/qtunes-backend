package router

import (
	"net/http"

	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo/v4"
)

// getUsers GET /users
func getUsersHandler(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := model.GetUsers(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, users)
}
