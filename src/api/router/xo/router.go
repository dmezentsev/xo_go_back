package xo

import (
	"api/app"
	"api/router"
	"api/router/protocol"
	"errors"
	"github.com/labstack/echo"
	"net/http"
)

type RouterContext struct {
	router.Context
}

func (r *RouterContext) NewGame(e echo.Context) error {
	request := new(protocol.CreateRoomRequest)
	if err := e.Bind(request); err != nil {
		return err
	}
	room, err := r.App.NewRoom(request.Name)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}
	game := NewGame(room)
	room.Meta = game
	return e.JSON(http.StatusOK, room)
}

func (r *RouterContext) NewPlayer(e echo.Context) error {
	room, err := r.App.GetRoom(app.UIDType(e.Param("roomUID")))
	if err != nil {
		return err
	}
	game := room.Meta.(*Game)
	if game == nil {
		return errors.New("nil game allocated")
	}

	if err := game.ValidateNewPlayer(); err != nil {
		return err
	}

	request := new(protocol.ParticipantRequest)
	if err := e.Bind(request); err != nil {
		return err
	}
	participant, err := room.NewParticipant(request.Name)
	if err != nil {
		return err
	}
	player, err := game.NewPlayer(participant)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, player.Participant)
}

func (r *RouterContext) NewWatcher(e echo.Context) error {
	room, err := r.App.GetRoom(app.UIDType(e.Param("roomUID")))
	if err != nil {
		return err
	}
	game := room.Meta.(*Game)
	if game == nil {
		return errors.New("nil game allocated")
	}

	request := new(protocol.ParticipantRequest)
	if err := e.Bind(request); err != nil {
		return err
	}
	participant, err := room.NewParticipant(request.Name)
	if err != nil {
		return err
	}
	watcher, err := game.NewWatcher(participant)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, watcher.Participant)
}

func (r *RouterContext) Connect(e echo.Context) error {
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
