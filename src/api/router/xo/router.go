package xo

import (
	"api/app"
	"api/router"
	"api/router/protocol"
	"github.com/labstack/echo"
	"net/http"
)

type RouterContext struct {
	router.Context
}

func (r *RouterContext) NewGame(e echo.Context) error {
	room, err := r.App.NewRoom()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}
	game := NewGame(room)
	room.Meta = game
	return e.JSON(http.StatusOK, protocol.RoomSerialize(room))
}

func (r *RouterContext) NewPlayer(e echo.Context) error {
	room, err := r.App.GetRoom(app.UIDType(e.Param("roomUID")))
	if err != nil {
		return err
	}

	game := room.Meta.(*Game)
	player, err := game.NewPlayer(room)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, protocol.ParticipantSerialize(player.Participant))
}

func (r *RouterContext) NewWatcher(e echo.Context) error {
	room, err := r.App.GetRoom(app.UIDType(e.Param("roomUID")))
	if err != nil {
		return err
	}

	game := room.Meta.(*Game)
	watcher, err := game.NewWatcher(room)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, protocol.ParticipantSerialize(watcher.Participant))
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
