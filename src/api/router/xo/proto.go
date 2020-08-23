package xo

import "api/app"

const BoardStateMessageType = "xo_board_state"

func BuildBoardState(b *Board) app.Message {
	return app.Message{
		Type: BoardStateMessageType,
		Payload: b.Fields,
	}
}
