package xo

import (
	"net/http"

	"github.com/labstack/echo"

	"api/router/protocol"
)

func (r *RouterContext) NewRoom(e echo.Context) error {
	room, err := r.App.NewRoom("xoHandler")
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, protocol.RoomSerialize(room))
}
