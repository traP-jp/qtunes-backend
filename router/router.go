package router

import (
	"errors"
	"net/http"
	"os"

	sess "github.com/hackathon-21-spring-02/back-end/session"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	s            sess.Session
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
)

// パッケージの初期化
func init() {
	if clientID == "" {
		panic(errors.New("clientID should not be empty."))
	}
	if clientSecret == "" {
		panic(errors.New("clientSecret should not be empty."))
	}
}

func SetRouting(sess sess.Session) {
	s = sess

	e := echo.New()

	e.Use(session.Middleware(sess.Store()))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/openapi", "docs/swagger")

	api := e.Group("/api")
	{
		apiPing := api.Group("/ping")
		{
			apiPing.GET("", func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusOK, "pong!")
			})
		}

		// ユーザー情報
		apiUsers := api.Group("/users")
		{
			apiUsers.GET("", getUsers , userAuthMiddleware)
		}

		apiFiles := api.Group("/files")
		{
			apiFiles.GET("",getFilesHandler,userAuthMiddleware)
		}
		// OAuth関連
		apiOAuth := api.Group("/oauth")
		{
			apiOAuth.GET("/callback", CallbackHandler)
			apiOAuth.GET("/generate/code", PostGenerateCodeHandler)
			apiOAuth.POST("/logout", PostLogoutHandler, userAuthMiddleware)
		}
	}

	err := e.Start(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}

