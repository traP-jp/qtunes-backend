package router

import (
	"fmt"
	"net/http"

	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo-contrib/session"
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

func getUserHandler(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("userID")

	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session:%w", err).Error())
	}
	accessToken := sess.Values["accessToken"].(string)
	res, err := model.GetUser(ctx, accessToken, userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get file: %w", err).Error())
	}
	return echo.NewHTTPError(http.StatusOK, res)
}
