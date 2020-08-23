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
	game := NewGame()
	room, err := r.App.NewRoom(game)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}
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
