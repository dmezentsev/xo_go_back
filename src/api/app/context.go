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
	rooms     []*RoomContext
	roomIndex map[UIDType]int16
	roomMutex sync.RWMutex
	Cancel    context.CancelFunc
	Mode      string
}

func NewApp(mode string) *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{
		context:   ctx,
		Cancel:    cancel,
		Mode:      mode,
		rooms:     make([]*RoomContext, 0),
		roomIndex: make(map[UIDType]int16),
	}
}
