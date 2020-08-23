package xo

import (
	"api/app"
	"sync"
)

type Board struct {
	State [][]SignType `json:"board"`
	mux sync.RWMutex
	BoardChanges chan event
}

func (b *Board) Move(m move) error {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.State[m.X][m.Y] = m.Sign
	b.BoardChanges <- event{Type: BoardChangesType}
	return nil
}

type Player struct {
	Participant *app.Participant
	Sign SignType
}
