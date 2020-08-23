package app

import (
	"context"
	"sync"
)

type UIDType string

type Context struct {
	context   context.Context
	rooms     map[UIDType]*RoomContext
	roomMutex sync.RWMutex
	Cancel    context.CancelFunc
}

func NewApp() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{context: ctx, Cancel: cancel, rooms: make(map[UIDType]*RoomContext)}
}
