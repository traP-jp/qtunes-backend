package router

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// userAuthMiddleware 本番用のAPIにアクセスしたユーザーを認証するミドルウェア
func userAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			return errSessionNotFound(err)
		}

		accessToken := sess.Values["accessToken"]
		if accessToken == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		c.Set("accessToken", accessToken)

		return next(c)
	}
}
