package router

import (
	"fmt"
	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

func getComposersHandler(c echo.Context) error {
	ctx:=c.Request().Context()
	sess,err:=session.Get("sessions",c)
	if err != nil {
		return c.String(http.StatusInternalServerError,fmt.Errorf("Failed in Getting Session:%w",err).Error())
	}

	accessToken:=sess.Values["accessToken"].(string)
	userID:=sess.Values["id"].(string)
	composers,err:=model.GetComposers(ctx,accessToken,userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,err)
	}

	return echo.NewHTTPError(http.StatusOK,composers)
}
