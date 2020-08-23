package router

import (
	"github.com/labstack/echo"
	"net/http"

	"api/app"
	"api/router/protocol"
)

func (r *Context) NewParticipant(e echo.Context) error {
	UID := e.Param("uid")
	room, err := r.App.GetRoom(app.UIDType(UID))
	if err != nil {
		return err
	}
	participant, err := room.NewParticipant()
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, protocol.ParticipantSerialize(participant))
}
