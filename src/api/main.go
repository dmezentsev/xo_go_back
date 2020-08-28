package main

import (
	"api/app"
	"api/router"
	"api/router/xo"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		fmt.Println(err)
		if err != nil {
			err = c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
		}
	}
	e.Logger.SetLevel(log.DEBUG)
	mode, exists := os.LookupEnv("MOD")
	if !exists {
		mode = app.ProdMode
	}
	appCtx := app.NewApp(mode)
	Router := router.Context{App: appCtx, Upgrader: app.Upgrader}
	xoRouter := xo.RouterContext{Context: Router}

	e.GET("/room", Router.GetRoomList)
	e.GET("/room/:uid", Router.GetRoom)
	e.DELETE("/room/:uid", Router.DeleteRoom)

	e.POST("/xo", xoRouter.NewGame)
	e.POST("/xo/:roomUID/player", xoRouter.NewPlayer)
	e.POST("/xo/:roomUID/watcher", xoRouter.NewWatcher)
	e.GET("/xo/:roomUID/connect/:participantUID", xoRouter.Connect)

	e.Logger.Debug(e.Start(":1323"))
}
