package main

import (
	"api/app"
	"api/router"
	"api/router/xo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	appCtx := app.NewApp()
	Router := router.Context{App: appCtx, Upgrader: app.Upgrader}
	xoRouter := xo.RouterContext{Context: Router}

	e.GET("/room/:uid", Router.GetRoom)
	e.POST("/room/:uid/participant", Router.NewParticipant)

	e.POST("/room/xo", xoRouter.NewRoom)
	e.POST("/room/xo/:roomUID/player", xoRouter.NewPlayer)
	e.GET("/room/xo/:roomUID/player/:participantUID/connect", xoRouter.ConnectPlayer)

	e.Logger.Debug(e.Start(":1323"))
}
