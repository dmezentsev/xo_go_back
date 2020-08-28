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
	Name              string       `json:"name"`
	UID               UIDType      `json:"uid"`
	Room              *RoomContext `json:"-"`
	Bus               *bus.Bus     `json:"-"`
	Absorber          chan Message `json:"-"`
	connectionEmitter bus.Emitter
	Connected         bool `json:"connected"`
	lastModified      time.Time
	Meta              interface{} `json:"meta"`
}

func (p *Participant) String() string {
	return fmt.Sprintf("<Participant: %s>", p.UID)
}

func (room *RoomContext) GetParticipant(participantUID UIDType) (*Participant, error) {
	room.participantMutex.RLock()
	defer room.participantMutex.RUnlock()
	participantIdx, ok := room.participantIndex[participantUID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("participant with UID: '%s' doesn't exists", participantUID))
	}
	return room.Participants[participantIdx], nil
}

func (room *RoomContext) GetParticipants() []*Participant {
	room.participantMutex.RLock()
	defer room.participantMutex.RUnlock()
	return room.Participants
}

func (room *RoomContext) NewParticipant(name string) (*Participant, error) {
	room.participantMutex.Lock()
	defer room.participantMutex.Unlock()
	participantUID := room.generateParticipantUID()
	participant := &Participant{
		Name:         name,
		UID:          participantUID,
		Room:         room,
		lastModified: time.Now(),
		Absorber:     make(chan Message),
		Connected:    false,
	}
	participant.Bus = bus.NewBus(participant.String())
	room.Participants = append(room.Participants, participant)
	room.participantIndex[participantUID] = int32(len(room.Participants) - 1)
	return participant, nil
}

func (room *RoomContext) DeleteParticipant(UID UIDType) error {
	room.participantMutex.Lock()
	defer room.participantMutex.Unlock()
	participantIdx, ok := room.participantIndex[UID]
	if !ok {
		return errors.New(fmt.Sprintf("participant with UID '%s' doesn't exists", UID))
	}
	room.Cancel()
	room.Participants = append(room.Participants[:participantIdx], room.Participants[participantIdx+1:]...)
	delete(room.participantIndex, UID)
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
		if _, ok := room.participantIndex[participantUid]; !ok {
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
