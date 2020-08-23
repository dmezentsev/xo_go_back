package xo

import (
	"api/app"
	"api/bus"
	"encoding/json"
	"errors"
	"fmt"
)

type Player struct {
	*app.Participant
	Sign SignType
	MoveEmitter bus.Emitter
}

const MoveEventType = bus.EventType("xo_move")

type MoveEvent struct {
	bus.Event
	X int `json:"x"`
	Y int `json:"y"`
}

func (g *Game) NewPlayer(room *app.RoomContext) (*Player, error) {
	g.mux.Lock()
	defer g.mux.Unlock()
	if len(g.players) >= 2 {
		return nil, errors.New("must be only two players")
	}
	participant, err := room.NewParticipant()
	if err != nil {
		return nil, err
	}
	var sign SignType
	if len(g.players) == 0 {
		sign = XSign
	} else {
		sign = OSign
	}
	player := &Player{
		Participant: participant,
		Sign: sign,
	}
	g.players = append(g.players, player)
	player.MoveEmitter = g.Bus.NewEmitter(MoveEventType, player)
	participant.Bus.NewCallback(app.MessageAcceptedEventType, player.onMessageReceive, nil)
	participant.Bus.NewCallback(app.ConnectEventType, player.onUserConnect, g.Board)
	g.Bus.NewCallback(BoardChangesEventType, player.onBoardChanged, nil)
	g.Bus.NewCallback(EndGameEventType, player.onEndGame, nil)
	return player, nil
}

func (player *Player) onBoardChanged(args bus.CallbackArgs) error {
	player.Absorber <- BuildBoardStateMessage(args.Initiator.(*Board))
	return nil
}

func (player *Player) onUserConnect(args bus.CallbackArgs) error {
	player.Absorber <- BuildBoardStateMessage(args.Meta.(*Board))
	return nil
}

func (player *Player) onMessageReceive(args bus.CallbackArgs) error {
	msg := bus.Event{}
	if err := json.Unmarshal(args.Event.GetPayload().([]byte), &msg); err != nil {
		return err
	}
	switch msg.Type {
	case MoveEventType:
		move := MoveEvent{}
		if err := json.Unmarshal(args.Event.GetPayload().([]byte), &move); err != nil {
			return err
		}
		player.MoveEmitter.Emitter <- move
	default:
		return errors.New(fmt.Sprintf(`unknown event type "%s"`, msg.Type))
	}
	return nil
}

func (player *Player) onEndGame(args bus.CallbackArgs) error {
	wonSign := args.Event.GetPayload().(SignType)
	var result GameResultType
	if wonSign == player.Sign {
		result = WinResult
	} else if wonSign == EmptySign {
		result = DrawResult
	} else {
		result = LooseResult
	}
	player.Absorber <- BuildPlayerEndGameMessage(result)
	return nil
}
