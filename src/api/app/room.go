package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type RoomContext struct {
	UID              UIDType
	App              *Context
	Meta          interface{}
	Bus *Bus
	context          context.Context
	Cancel           context.CancelFunc
	lastModified     time.Time
	Participants     map[UIDType]*Participant
	participantMutex sync.RWMutex
}

func (app *Context) GetRoom(roomUID UIDType) (*RoomContext, error) {
	app.roomMutex.RLock()
	defer app.roomMutex.RUnlock()
	room, ok := app.rooms[roomUID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("room with UID: '%s' doesn't exists", roomUID))
	}
	return room, nil
}

func (app *Context) NewRoom(meta interface{}) (*RoomContext, error) {
	app.roomMutex.Lock()
	defer app.roomMutex.Unlock()
	roomUID := app.generateRoomUID()
	ctx, cancel := context.WithCancel(context.Background())
	bus := NewBus()
	c := func() {
		bus.Cancel()
		cancel()
	}
	participants := make(map[UIDType]*Participant)
	room := &RoomContext{
		UID:          roomUID,
		App:          app,
		Meta:      meta,
		Bus: NewBus(),
		lastModified: time.Now(),
		context:      ctx,
		Cancel:       c,
		Participants: participants,
	}
	app.rooms[roomUID] = room
	return room, nil
}

func (app *Context) DeleteRoom(roomUID UIDType) error {
	app.roomMutex.Lock()
	defer app.roomMutex.Unlock()
	room, ok := app.rooms[roomUID]
	if !ok {
		return errors.New(fmt.Sprintf("roomUID '%s' doesn't exists", roomUID))
	}
	room.Cancel()
	delete(app.rooms, roomUID)
	return nil
}

func (app *Context) generateRoomUID() UIDType {
	var roomUid UIDType
	for {
		roomUid = UIDType(uuid.New().String())
		if _, ok := app.rooms[roomUid]; !ok {
			break
		}
	}
	return "9ae53419-a0c4-4c53-ab06-14c1dcb5808b"
}
