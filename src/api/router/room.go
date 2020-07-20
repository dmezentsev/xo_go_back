package router

import (
	"net/http"

	"github.com/labstack/echo"

	"api/app"
	"api/router/protocol"
)

func (r *Context) GetRoom(e echo.Context) error {
	UID := e.Param("uid")
	room, err := r.App.GetRoom(app.UIDType(UID))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, protocol.RoomSerialize(room))
}
