package router

import (
	"net/http"

	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// getComposersHandler GET /composers
func getComposersHandler(c echo.Context) error {
	ctx := c.Request().Context()
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	accessToken := sess.Values["accessToken"].(string)
	composers, err := model.GetComposers(ctx, accessToken)
	if err != nil {
		return handleError(err)
	}

	return echo.NewHTTPError(http.StatusOK, composers)
}

// getComposerHandler GET /composers/:composerID
func getComposerHandler(c echo.Context) error {
	ctx := c.Request().Context()
	composerID := c.Param("composerID")
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	accessToken := sess.Values["accessToken"].(string)
	res, err := model.GetComposer(ctx, accessToken, composerID)
	if err != nil {
		return handleError(err)
	}

	return echo.NewHTTPError(http.StatusOK, res)
}

// getComposerFilesHandler GET /composers/:composerID/files
func getComposerFilesHandler(c echo.Context) error {
	ctx := c.Request().Context()
	composerID := c.Param("composerID")
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	accessToken := sess.Values["accessToken"].(string)
	userID := sess.Values["id"].(string)
	res, err := model.GetComposerFiles(ctx, accessToken, composerID, userID)
	if err != nil {
		return handleError(err)
	}

	return echo.NewHTTPError(http.StatusOK, res)
}

// getComposerByNameHandler GET /composers/name/:composerName
func getComposerByNameHandler(c echo.Context) error {
	ctx := c.Request().Context()
	name := c.Param("composerName")
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	accessToken := sess.Values["accessToken"].(string)
	res, err := model.GetComposerByName(ctx, accessToken, name)
	if err != nil {
		return handleError(err)
	}

	return echo.NewHTTPError(http.StatusOK, res)
}
