package xo

import (
	"api/app"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
)

type Player struct {
	*app.Participant
	Sign SignType
	MoveEmitter app.Emitter
}

const MoveEventType = app.EventType("xo_move")

type MoveEvent struct {
	app.Event
	X int `json:"x"`
	Y int `json:"y"`
}

type EndGameResultType string
//const WinResult = EndGameResultType("xo_win")
//const LooseResult = EndGameResultType("xo_loose")
//const DrawResult = EndGameResultType("xo_draw")

type EndGameEvent struct {
	Result EndGameResultType `json:"result"`
}

func (g *Game) NewPlayer(room *app.RoomContext) (*Player, error) {
	participant, err := room.NewParticipant()
	if err != nil {
		return nil, err
	}
	player := &Player{
		Participant: participant,
		Sign: XSign, // TODO: define as rule
	}
	player.MoveEmitter = g.Bus.NewEmitter(MoveEventType, player)
	participant.Bus.NewCallback(app.MessageAcceptedEvent, player.onMessageReceive, nil)
	g.Bus.NewCallback(BoardChangesEventType, player.onBoardChanged, nil)
	g.Bus.NewCallback(EndGameEventType, player.onBoardChanged, nil)
	return player, nil
}

func (player *Player) onMessageReceive(args app.CallbackArgs) {
	fmt.Printf("%+v\n", args)
	msg := app.Event{}
	if err := json.Unmarshal(args.Event.GetPayload().([]byte), &msg); err != nil {
		return
	}
	switch msg.Type {
	case MoveEventType:
		move := MoveEvent{}
		if err := json.Unmarshal(args.Event.GetPayload().([]byte), &move); err != nil {
			return
		}
		player.MoveEmitter.Emitter <- move
	}
}

func (player *Player) onBoardChanged(args app.CallbackArgs) {
	fmt.Printf("%+v\n", args)
	player.Participant.Absorber <- BuildBoardState(args.Initiator.(*Game).Board)
}

func (player *Player) onEndGame(args app.CallbackArgs) {
	fmt.Printf("%+v\n", args)
	fmt.Println("Calculate won player + Emit message to participant Absorber?")
}

func (r *RouterContext) ConnectPlayer(e echo.Context) error {
	var err error

	defer func() {
		if err != nil {
			e.Logger().Error(err)
		}
	}()

	e.Logger().Debug("Initialize client start")
	room, err := r.App.GetRoom(app.UIDType(e.Param("roomUID")))
	if err != nil {
		e.Logger().Error(err)
		return err
	}
	participant, err := room.GetParticipant(app.UIDType(e.Param("participantUID")))
	if err != nil {
		e.Logger().Error(err)
		return err
	}

	transport, err := app.NewWsTransport(e)
	if err != nil {
		e.Logger().Error(err)
		return err
	}

	if err = participant.Connect(transport); err != nil {
		e.Logger().Error(err)
		return err
	}

	return nil
}
