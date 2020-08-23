package xo

import (
	"api/app"
	"fmt"
	"sync"
)

const BoardChangesEventType = app.EventType("xo_board_changes")
const EndGameEventType = app.EventType("xo_end_game")

type boardChangesEvent struct {
	app.Event
}

type endGameEvent struct {
	app.Event
}

type Board struct {
	Fields         [][]SignType `json:"board"`
	changesEmitter app.Emitter
	mux            sync.RWMutex
}

type Game struct {
	Board *Board
	Bus *app.Bus
	endGameEmitter app.Emitter
}

func NewBoard(rows int, cols int) *Board {
	b := Board{}
	b.mux.Lock()
	defer b.mux.Unlock()
	for r := 0; r < rows; r++ {
		b.Fields = append(b.Fields, make([]SignType, cols))
	}
	return &b
}

func (b *Board) Move(sign SignType, x int, y int) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Fields[x][y] = sign
	b.changesEmitter.Emitter <- boardChangesEvent{}
}

func (b *Board) IsEndGame() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	for _, rows := range b.Fields {
		for _, y := range rows {
			if y == EmptySign {
				return false
			}
		}
	}
	return true
}

func NewGame() *Game {
	bus := app.NewBus()
	board := NewBoard(3, 3)
	g := &Game{
		Board: board,
		Bus: bus,
	}
	g.Board.changesEmitter = bus.NewEmitter(BoardChangesEventType, board)
	g.Board.changesEmitter = bus.NewEmitter(EndGameEventType, g)
	bus.NewCallback(MoveEventType, g.onMove, nil)
	return g
}

func (g *Game) onMove(args app.CallbackArgs) {
	fmt.Println("Move rules apply + Board.Move")
	player := args.Initiator.(*Player)
	move := args.Event.(MoveEvent)
	g.Board.Move(player.Sign, move.X, move.Y)
	if g.Board.IsEndGame() {
		g.endGameEmitter.Emitter <- endGameEvent{}
	}
}
