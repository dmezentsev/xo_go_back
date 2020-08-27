package app

import (
	"api/bus"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type RoomContext struct {
	Name             string `json:"name"`
	UID              UIDType `json:"uid"`
	App              *Context `json:"-"`
	Meta             interface{} `json:"meta"`
	Bus              *bus.Bus `json:"-"`
	context          context.Context
	Cancel           context.CancelFunc `json:"-"`
	lastModified     time.Time
	Participants     []*Participant `json:"participant"`
	participantIndex map[UIDType]int32
	participantMutex sync.RWMutex
}

func (app *Context) GetRoomList() ([]*RoomContext, error) {
	app.roomMutex.RLock()
	defer app.roomMutex.RUnlock()
	return app.rooms, nil
}

func (app *Context) GetRoom(roomUID UIDType) (*RoomContext, error) {
	app.roomMutex.RLock()
	defer app.roomMutex.RUnlock()
	roomIdx, ok := app.roomIndex[roomUID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("room with UID: '%s' doesn't exists", roomUID))
	}
	room := app.rooms[roomIdx]
	return room, nil
}

func (app *Context) NewRoom(name string) (*RoomContext, error) {
	app.roomMutex.Lock()
	defer app.roomMutex.Unlock()
	roomUID := app.generateRoomUID()
	ctx, cancel := context.WithCancel(context.Background())
	room := &RoomContext{
		Name:             name,
		UID:              roomUID,
		App:              app,
		lastModified:     time.Now(),
		context:          ctx,
		Participants:     make([]*Participant, 0),
		participantIndex: make(map[UIDType]int32),
	}
	room.Bus = bus.NewBus(room.String())
	room.Cancel = func() {
		room.Bus.Cancel()
		cancel()
	}
	app.rooms = append(app.rooms, room)
	app.roomIndex[roomUID] = int16(len(app.rooms) - 1)
	return room, nil
}

func (room *RoomContext) String() string {
	return fmt.Sprintf("<Room %s %+v>", room.UID, room.Meta)
}

func (app *Context) DeleteRoom(roomUID UIDType) error {
	app.roomMutex.Lock()
	defer app.roomMutex.Unlock()
	roomIdx, ok := app.roomIndex[roomUID]
	if !ok {
		return errors.New(fmt.Sprintf("roomUID '%s' doesn't exists", roomUID))
	}
	room := app.rooms[roomIdx]
	room.Cancel()
	delete(app.roomIndex, roomUID)
	app.rooms = append(app.rooms[:roomIdx], app.rooms[roomIdx+1:]...)
	return nil
}

func (app *Context) generateRoomUID() UIDType {
	var roomUid UIDType
	for {
		roomUid = UIDType(uuid.New().String())
		if _, ok := app.roomIndex[roomUid]; !ok {
			break
		}
	}
	if app.Mode == DebugMode {
		return "9ae53419-a0c4-4c53-ab06-14c1dcb5808b"
	}
	return roomUid
}
