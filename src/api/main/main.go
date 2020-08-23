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

	e.POST("/xo", xoRouter.NewGame)
	e.POST("/xo/:roomUID/player", xoRouter.NewPlayer)
	e.POST("/xo/:roomUID/watcher", xoRouter.NewWatcher)
	e.GET("/xo/:roomUID/connect/:participantUID", xoRouter.Connect)

	e.Logger.Debug(e.Start(":1323"))
}
