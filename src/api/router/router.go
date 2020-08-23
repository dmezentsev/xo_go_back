package router

import (
	"api/app"
	"github.com/gorilla/websocket"
)

type Context struct {
	App      *app.Context
	Upgrader websocket.Upgrader
}
