package app

import (
	"context"
	"sync"
)

type UIDType string

const DebugMode = "debug"
const TestingMode = "test"
const ProdMode = "prod"

type Context struct {
	context   context.Context
	rooms     map[UIDType]*RoomContext
	roomMutex sync.RWMutex
	Cancel    context.CancelFunc
	Mode string
}

func NewApp(mode string) *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{
		context: ctx,
		Cancel: cancel,
		Mode: mode,
		rooms: make(map[UIDType]*RoomContext),
	}
}
