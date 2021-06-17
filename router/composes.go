package router

import (
	"fmt"
	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo/v4"
)

func composersHandler(c echo.Context) error {
	composers,err:=model.GetComposers()
	if err != nil {
		return err
	}
	fmt.Println(composers)
	return err
}