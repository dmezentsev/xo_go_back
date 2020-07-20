package main

import (
	"api/router"
	"api/router/xo"
	"github.com/labstack/echo"

	"api/app"
)

func main() {
	e := echo.New()
	appCtx := app.NewApp()
	Router := router.Context{App: appCtx}
	xoRouter := &xo.RouterContext{Context: Router}

	e.POST("/room/xo", xoRouter.NewRoom)
	e.GET("/room/:uid", Router.GetRoom)

	e.Logger.Debug(e.Start(":1323"))
}
