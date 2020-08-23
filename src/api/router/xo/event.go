package xo

type IEvent interface {
	GetType() string
}

type event struct {
	IEvent
	Initiator interface{}  `json:"-"`
	Type string  `json:"t"`
}

type Connect struct {
	event
}

type Disconnect struct {
	event
}

type IMove interface {
	IEvent
}

type move struct {
	event
	X         int `json:"x"`
	Y         int `json:"y"`
	Sign SignType
}

const BoardChangesType = "board_changes"
