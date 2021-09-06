package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type editFavoriteRequest struct {
	Favorite *bool `json:"favorite"`
}

// getFilesHandler GET /files
func getFilesHandler(c echo.Context) error {
	ctx := c.Request().Context()
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	userID := sess.Values["id"].(string)
	files, err := model.GetFiles(ctx, userID)
	if err != nil {
		return generateEchoError(err)
	}

	return echo.NewHTTPError(http.StatusOK, files)
}

// getRandomFileHandler GET /files/random
func getRandomFileHandler(c echo.Context) error {
	ctx := c.Request().Context()
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	userID := sess.Values["id"].(string)
	file, err := model.GetRandomFile(ctx, userID)
	if err != nil {
		return generateEchoError(err)
	}

	return echo.NewHTTPError(http.StatusOK, file)
}

// getFileHandler GET /files/:fileID
func getFileHandler(c echo.Context) error {
	ctx := c.Request().Context()
	fileID := c.Param("fileID")
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	userID := sess.Values["id"].(string)
	file, err := model.GetFile(ctx, userID, fileID)
	if err != nil {
		return generateEchoError(err)
	}

	return echo.NewHTTPError(http.StatusOK, file)
}

// getFileDownloadHandler GET /files/:fileID/download
func getFileDownloadHandler(c echo.Context) error {
	ctx := c.Request().Context()
	fileID := c.Param("fileID")

	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	accessToken := sess.Values["accessToken"].(string)

	file, res, err := model.GetFileDownload(ctx, fileID, accessToken)
	if err != nil {
		return generateEchoError(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return generateEchoError(err)
	}

	c.Response().Header().Set(echo.HeaderContentType, res.Header.Get("Content-Type"))
	c.Response().Header().Set("Cache-Control", "private, max-age=31536000") // 1年間キャッシュ
	http.ServeContent(c.Response(), c.Request(), info.Name(), info.ModTime(), file)
	return echo.NewHTTPError(http.StatusOK)
}

// putFileFavoriteHandler PUT /files/:fileID/favorite
func putFileFavoriteHandler(c echo.Context) error {
	ctx := c.Request().Context()
	fileID := c.Param("fileID")
	fav := editFavoriteRequest{}
	err := c.Bind(&fav)
	if err != nil {
		return errBind(err)
	}
	if fav.Favorite == nil {
		return errBind(errors.New("invalid type"))
	}
	log.Printf("%+v", fav)

	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	userID := sess.Values["id"].(string)
	err = model.ToggleFileFavorite(ctx, userID, fileID, *fav.Favorite)
	if err != nil {
		return generateEchoError(err)
	}

	return echo.NewHTTPError(http.StatusOK)
}

// getFileFromTitleHandler GET /files/title/:title
func getFileFromTitleHandler(c echo.Context) error {
	ctx := c.Request().Context()
	title := c.Param("title")
	_, err := session.Get("sessions", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session: %w", err))
	}
	file, err := model.FindFileFromTitle(ctx,title)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, file)
}

// getFileFromComposerNameHandler GET /files/composer/:composerName
func getFileFromComposerNameHandler(c echo.Context) error {
	ctx := c.Request().Context()
	composerName := c.Param("composerName")
	_, err := session.Get("sessions", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session: %w", err))
	}
	file, err := model.FindFileFromComposerName(ctx,composerName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, file)
}
