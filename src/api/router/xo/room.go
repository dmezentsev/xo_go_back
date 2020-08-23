package xo

import (
	"api/router/protocol"
	"net/http"

	"github.com/labstack/echo"
)

func (r *RouterContext) NewRoom(e echo.Context) error {
	room, err := r.App.NewRoom("xoHandler")
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}
	return e.JSON(http.StatusOK, protocol.RoomSerialize(room))
}
