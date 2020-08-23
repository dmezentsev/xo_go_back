package xo

import (
	"api/app"
	"api/bus"
	"errors"
	"fmt"
	"sync"
)

const BoardChangesEventType = bus.EventType("xo_board_changes")
const EndGameEventType = bus.EventType("xo_end_game")

type Board struct {
	Fields          [][]SignType `json:"board"`
	changesEmitter  bus.Emitter
	endGameEmitter  bus.Emitter
	mux             sync.RWMutex
	currentSignMove SignType
	rows            int
	cols            int
	complete        bool
}

type Game struct {
	Room    *app.RoomContext
	Board   *Board
	Bus     *bus.Bus
	mux     sync.RWMutex
	players []*Player
}

func NewGame(room *app.RoomContext) *Game {
	board := NewBoard(3, 3)
	g := &Game{
		Room:    room,
		Board:   board,
		players: make([]*Player, 0),
	}
	g.Bus = bus.NewBus(g.String())
	g.Board.changesEmitter = g.Bus.NewEmitter(BoardChangesEventType, board, nil)
	g.Board.endGameEmitter = g.Bus.NewEmitter(EndGameEventType, board, nil)
	g.Bus.NewCallback(MoveEventType, g.onMove, nil)
	return g
}

func (g *Game) String() string {
	return fmt.Sprintf("<Game XO %s>", g.Room.UID)
}

func (g *Game) onMove(args bus.CallbackArgs) error {
	player := args.Initiator.(*Player)
	move := args.Event.(MoveEvent)
	return g.Board.Move(player.Sign, move.X, move.Y)
}

func NewBoard(rows int, cols int) *Board {
	b := Board{currentSignMove: XSign, rows: rows, cols: cols}
	b.mux.Lock()
	defer b.mux.Unlock()
	for r := 0; r < rows; r++ {
		b.Fields = append(b.Fields, make([]SignType, 0))
		for c := 0; c < cols; c++ {
			b.Fields[r] = append(b.Fields[r], EmptySign)
		}
	}
	return &b
}

func (b *Board) Move(sign SignType, x int, y int) error {
	if x >= b.rows || y >= b.cols || x < 0 || y < 0 {
		return errors.New(fmt.Sprintf(`unbound {x:y} values {%d:%d}`, x, y))
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	if b.complete {
		return errors.New("game is finish")
	}
	if b.Fields[x][y] != EmptySign {
		return errors.New("field is occupied")
	}
	if b.currentSignMove != sign {
		return errors.New("waiting opponent moving")
	}
	b.Fields[x][y] = sign
	if sign == XSign {
		b.currentSignMove = OSign
	} else {
		b.currentSignMove = XSign
	}
	b.changesEmitter.Emitter <- bus.Event{}
	winner := b.GetWinner()
	if b.AllFieldFills() || winner != EmptySign {
		b.complete = true
		b.endGameEmitter.Emitter <- bus.Event{Payload: winner}
	}
	return nil
}

func (b *Board) AllFieldFills() bool {
	for _, rows := range b.Fields {
		for _, y := range rows {
			if y == EmptySign {
				return false
			}
		}
	}
	return true
}

func (b *Board) GetWinner() SignType {
	winner := EmptySign
	for _, row := range b.Fields {
		if winner = winByRow(row); winner != EmptySign {
			return winner
		}
	}
	for y := 0; y < b.cols; y++ {
		row := make([]SignType, b.rows)
		for x := 0; x < b.rows; x++ {
			row[x] = b.Fields[x][y]
		}
		if winner = winByRow(row); winner != EmptySign {
			return winner
		}
	}
	if b.rows == b.cols {
		for _, mod := range []int{0, b.rows - 1} {
			row := make([]SignType, b.rows)
			for i := 0; i < b.rows; i++ {
				x := i
				if mod != 0 {
					x = mod - i
				}
				row[i] = b.Fields[x][i]
			}
			if winner = winByRow(row); winner != EmptySign {
				return winner
			}
		}
	}
	return EmptySign
}

func winByRow(row []SignType) SignType {
	possibleWinner := EmptySign
	for _, v := range row {
		if v == EmptySign {
			return EmptySign
		}
		if possibleWinner == EmptySign {
			possibleWinner = v
		}
		if possibleWinner != v {
			return EmptySign
		}
	}
	return possibleWinner
}
