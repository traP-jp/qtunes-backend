package router

import (
	"fmt"
	"net/http"

	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type editFavoriteRequest struct {
	Favorite bool
}

// getFilesHandler GET /files
func getFilesHandler(c echo.Context) error {
	ctx := c.Request().Context()
	sess, err := session.Get("sessions", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session:%w", err))
	}
	userID := sess.Values["id"].(string)
	files, err := model.GetFiles(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, files)
}

// getRandomFileHandler GET /files/random
func getRandomFileHandler(c echo.Context) error {
	ctx := c.Request().Context()
	sess, err := session.Get("sessions", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session:%w", err))
	}
	userID := sess.Values["id"].(string)
	file, err := model.GetRandomFile(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, file)
}

// getFileHandler GET /files/:fileID
func getFileHandler(c echo.Context) error {
	ctx := c.Request().Context()
	fileID := c.Param("fileID")
	sess, err := session.Get("sessions", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session: %w", err))
	}
	userID := sess.Values["id"].(string)
	file, err := model.GetFile(ctx, userID, fileID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, file)
}

// getFileDownloadHandler GET /files/:fileID/download
func getFileDownloadHandler(c echo.Context) error {
	ctx := c.Request().Context()
	fileID := c.Param("fileID")

	sess, err := session.Get("sessions", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get session: %w", err))
	}
	accessToken := sess.Values["accessToken"].(string)

	file, res, err := model.GetFileDownload(ctx, fileID, accessToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get file: %w", err))
	}

	return c.Stream(http.StatusOK, res.Header.Get("Content-Type"), file)
}

// putFileFavoriteHandler PUT /files/:fileID/favorite
func putFileFavoriteHandler(c echo.Context) error {
	ctx := c.Request().Context()
	fileID := c.Param("fileID")
	fav := editFavoriteRequest{}
	err := c.Bind(&fav)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Failed to bind request: %w", err))
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session: %w", err))
	}
	userID := sess.Values["id"].(string)
	err = model.ToggleFileFavorite(ctx, userID, fileID, fav.Favorite)
	if err == model.DBErrs["NoChange"] {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return echo.NewHTTPError(http.StatusOK)
}
