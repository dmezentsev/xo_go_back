package router

import (
	"net/http"

	"github.com/labstack/echo"

	"api/app"
)

func (r *Context) GetRoomList(e echo.Context) error {
	rooms, err := r.App.GetRoomList()
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, rooms)
}

func (r *Context) GetRoom(e echo.Context) error {
	UID := e.Param("uid")
	room, err := r.App.GetRoom(app.UIDType(UID))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, room)
}

func (r *Context) DeleteRoom(e echo.Context) error {
	UID := e.Param("uid")
	if err := r.App.DeleteRoom(app.UIDType(UID)); err != nil {
		return err
	}
	return e.NoContent(http.StatusOK)
}
