package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type UIDType string
type ParticipantStatus string

const ParticipantStatusConnected = "connected"
const ParticipantStatusDisconnected = "disconnected"

type Context struct {
	context   context.Context
	rooms     map[UIDType]*RoomContext
	roomMutex sync.RWMutex
	Cancel    context.CancelFunc
}

type RoomContext struct {
	UID              UIDType
	App              *Context
	Handler          string // TODO: must be func()
	context          context.Context
	Cancel           context.CancelFunc
	lastModified     time.Time
	Participants     map[UIDType]*Participant
	participantMutex sync.RWMutex
}

type Participant struct {
	UID          UIDType
	Status       ParticipantStatus
	Room         *RoomContext
	Emitter      chan []byte
	Absorber     chan []byte
	lastModified time.Time
}

func NewApp() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{context: ctx, Cancel: cancel}
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

func (app *Context) NewRoom(handler string) (*RoomContext, error) {
	app.roomMutex.Lock()
	defer app.roomMutex.Unlock()
	roomUID := app.generateRoomUID()
	ctx, cancel := context.WithCancel(context.Background())
	c := func () {
		// TODO: some notifications
		cancel()
	}
	room := &RoomContext{UID: roomUID, App: app, Handler: handler, lastModified: time.Now(), context: ctx, Cancel: c}
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
	return roomUid
}

func (room *RoomContext) GetParticipant(participantUID UIDType) (*Participant, error) {
	room.participantMutex.RLock()
	defer room.participantMutex.RUnlock()
	participant, ok := room.Participants[participantUID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("participant with UID: '%s' doesn't exists", participantUID))
	}
	return participant, nil
}

func (room *RoomContext) NewParticipant() (*Participant, error) {
	room.participantMutex.Lock()
	defer room.participantMutex.Unlock()
	participantUID := room.generateParticipantUID()
	participant := &Participant{UID: participantUID, Room: room, Status: ParticipantStatusDisconnected, lastModified: time.Now()}
	room.Participants[participantUID] = participant
	return participant, nil
}

func (room *RoomContext) DeleteParticipant(UID UIDType) error {
	room.participantMutex.Lock()
	defer room.participantMutex.Unlock()
	_, ok := room.Participants[UID]
	if !ok {
		return errors.New(fmt.Sprintf("participant with UID '%s' doesn't exists", UID))
	}
	room.Cancel()
	delete(room.Participants, UID)
	return nil
}

func (room *RoomContext) generateParticipantUID() UIDType {
	var participantUid UIDType
	for {
		participantUid = UIDType(uuid.New().String())
		if _, ok := room.Participants[participantUid]; !ok {
			break
		}
	}
	return participantUid
}
