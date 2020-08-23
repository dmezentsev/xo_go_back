package app

import "errors"

type XOHandler struct {
	UID     UIDType
	PlayerX *Participant
	PlayerO *Participant
}

func NewHandler(r *RoomContext) {

}

func (h *XOHandler) NewPlayer(p *Participant) error {
	if h.PlayerX == nil || h.PlayerX.UID == p.UID {
		h.PlayerX = p
		return nil
	}
	if h.PlayerO == nil || h.PlayerO.UID == p.UID {
		h.PlayerO = p
		return nil
	}
	return errors.New("must be only 2 players")
}

func (h *XOHandler) RemovePlayer(p *Participant) error {
	if h.PlayerX.UID == p.UID {
		h.PlayerX = nil
		return nil
	}
	if h.PlayerO == h.PlayerO {
		h.PlayerO = nil
		return nil
	}
	return errors.New("there is no player")
}
