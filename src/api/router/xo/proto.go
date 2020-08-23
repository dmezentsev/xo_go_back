package xo

import "api/app"

const BoardStateMessageType = app.MessageType("xo_board_state")
const GamePlayerResultMessageType = app.MessageType("xo_player_game_result")
const GameWatcherResultMessageType = app.MessageType("xo_watcher_game_result")
const MoveErrorMessageType = app.MessageType("xo_move_error")

func BuildBoardStateMessage(b *Board) app.Message {
	return app.Message{
		Type: BoardStateMessageType,
		Payload: b.Fields,
	}
}

type GameResultType string
const WinResult = GameResultType("xo_win")
const LooseResult = GameResultType("xo_loose")
const DrawResult = GameResultType("xo_draw")

func BuildPlayerEndGameMessage(result GameResultType) app.Message {
	return app.Message{
		Type:    GamePlayerResultMessageType,
		Payload: result,
	}
}

func BuildWatcherEndGameMessage(result SignType) app.Message {
	return app.Message{
		Type:    GameWatcherResultMessageType,
		Payload: result,
	}
}


func BuildErrorMoveMessage(err error) app.Message {
	return app.Message{
		Type: MoveErrorMessageType,
		Payload: err.Error(),
	}
}
