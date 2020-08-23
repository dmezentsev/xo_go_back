package app

import (
	"api/bus"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Participant struct {
	UID               UIDType
	Room              *RoomContext
	Bus               *bus.Bus
	Absorber          chan Message
	connectionEmitter bus.Emitter
	Connected         bool
	lastModified      time.Time
	Meta              interface{}
}

func (p *Participant) String() string {
	return fmt.Sprintf("<Participant: %s>", p.UID)
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
		Absorber:     make(chan Message),
		Connected:    false,
	}
	participant.Bus = bus.NewBus(participant.String())
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

const MessageAcceptedEventType = bus.EventType("message_accepted")
const ConnectEventType = bus.EventType("connect")

func (p *Participant) Connect(transport *Transport) error {
	if p.Connected {
		return errors.New("participant already connected")
	}
	p.Connected = true
	emitter := p.Bus.NewEmitter(MessageAcceptedEventType, p, nil)

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
				data, err := json.Marshal(msg)
				if err != nil {
					continue
				}
				transport.Writer <- data
			case <-transport.Context.Done():
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	go func() {
		defer func() {
			p.Bus.Cancel()
			transport.Cancel()
			p.Disconnect()
		}()

		for {
			select {
			case msg := <-transport.Reader:
				emitter.Emitter <- bus.Event{Payload: msg}
			case <-transport.Context.Done():
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	p.Bus.NewEmitter(ConnectEventType, p, nil).Emitter <- bus.Event{}
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
	if room.App.Mode == DebugMode {
		uids := []UIDType{
			"9ae53419-a0c4-4c53-ab06-14c1dcb5808b",
			"_9ae53419-a0c4-4c53-ab06-14c1dcb5808b",
			"~9ae53419-a0c4-4c53-ab06-14c1dcb5808b",
			"=9ae53419-a0c4-4c53-ab06-14c1dcb5808b",
		}
		return uids[len(room.Participants)%len(uids)]
	}
	return participantUid
}
