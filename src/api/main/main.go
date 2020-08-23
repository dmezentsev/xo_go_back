package main

import (
	"api/app"
	"api/router"
	"api/router/xo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"os"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	mode, exists := os.LookupEnv("MOD")
	if !exists {
		mode = app.ProdMode
	}
	appCtx := app.NewApp(mode)
	Router := router.Context{App: appCtx, Upgrader: app.Upgrader}
	xoRouter := xo.RouterContext{Context: Router}

	e.POST("/xo", xoRouter.NewGame)
	e.POST("/xo/:roomUID/player", xoRouter.NewPlayer)
	e.POST("/xo/:roomUID/watcher", xoRouter.NewWatcher)
	e.GET("/xo/:roomUID/connect/:participantUID", xoRouter.Connect)

	e.Logger.Debug(e.Start(":1323"))
}
