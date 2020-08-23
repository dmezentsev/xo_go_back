package xo

import (
	"github.com/labstack/echo"
	"net/http"

	"api/app"
	"api/router/protocol"
)

func (r *RouterContext) NewPlayer(e echo.Context) error {
	room, err := r.App.GetRoom(app.UIDType(e.Param("uid")))
	if err != nil {
		return err
	}
	participant, err := room.NewParticipant()
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, protocol.ParticipantSerialize(participant))
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
