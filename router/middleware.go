
package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo-contrib/session"
)

// userAuthMiddleware 本番用のAPIにアクセスしたユーザーを認証するミドルウェア
func userAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Failed to get session: %w", err))
		}

		accessToken := sess.Values["accessToken"]
		if accessToken == nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		c.Set("accessToken", accessToken)

		return next(c)
	}
}
