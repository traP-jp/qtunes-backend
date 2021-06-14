package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api")
	v1 := api.Group("/v1")
	{
		apiPing := v1.Group("/ping")
		{
			apiPing.GET("", func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusOK, "pong!")
			})
		}
	}

	err := e.Start(":1325")
	if err != nil {
		panic(err)
	}
}

