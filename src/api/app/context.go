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
	Room         *RoomContext
	Emitter      chan []byte
	Absorber     chan []byte
	Connected    bool
	lastModified time.Time
	Meta         interface{}
}

func NewApp() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{context: ctx, Cancel: cancel, rooms: make(map[UIDType]*RoomContext)}
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
	c := func() {
		// TODO: some notifications
		cancel()
	}
	participants := make(map[UIDType]*Participant)
	room := &RoomContext{
		UID:          roomUID,
		App:          app,
		Handler:      handler,
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

func (room *RoomContext) GetParticipant(participantUID UIDType) (*Participant, error) {
	room.participantMutex.RLock()
	defer room.participantMutex.RUnlock()
	participant, ok := room.Participants[participantUID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("participant with UID: '%s' doesn't exists", participantUID))
	}
	return participant, nil
}

func (room *RoomContext) GetParticipants() []*Participant {
	room.participantMutex.RLock()
	defer room.participantMutex.RUnlock()
	participants := make([]*Participant, len(room.Participants))
	i := 0
	for _, p := range room.Participants {
		participants[i] = p
		i++
	}
	return participants
}

func (room *RoomContext) NewParticipant() (*Participant, error) {
	room.participantMutex.Lock()
	defer room.participantMutex.Unlock()
	participantUID := room.generateParticipantUID()
	participant := &Participant{
		UID:          participantUID,
		Room:         room,
		lastModified: time.Now(),
		Emitter:      make(chan []byte),
		Absorber:     make(chan []byte),
		Connected:    false,
	}
	room.Participants[participantUID] = participant
	//go func() {
	//	for {
	//		fmt.Println("loop back")
	//		select {
	//		case msg := <-participant.Emitter:
	//			fmt.Printf("emmit message %s\n", msg)
	//			participant.Absorber <- msg
	//		default:
	//			time.Sleep(3 * time.Second)
	//		}
	//	}
	//}()
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

func (p *Participant) Connect(transport *Transport) error {
	if p.Connected {
		return errors.New("participant already connected")
	}
	p.Connected = true

	go func() {
		defer func() {
			transport.Cancel()
			p.Disconnect()
		}()

		for {
			if !p.Connected {
				return
			}
			select {
			case msg := <-p.Absorber:
				transport.Writer <- msg
			case <- transport.Context.Done():
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	go func() {
		defer func() {
			transport.Cancel()
			p.Disconnect()
		}()

		for {
			select {
			case msg := <-transport.Reader:
				p.Emitter <- msg
			case <- transport.Context.Done():
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	return nil
}

func (p *Participant) Disconnect() {
	p.Connected = false
}

func (room *RoomContext) generateParticipantUID() UIDType {
	var participantUid UIDType
	for {
		participantUid = UIDType(uuid.New().String())
		if _, ok := room.Participants[participantUid]; !ok {
			break
		}
	}
	return "9ae53419-a0c4-4c53-ab06-14c1dcb5808b"
}
