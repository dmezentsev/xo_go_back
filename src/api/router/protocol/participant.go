package protocol

import "api/app"

type ParticipantResponse struct {
	UID app.UIDType `json:"uid"`
}

func ParticipantSerialize(p *app.Participant) ParticipantResponse {
	return ParticipantResponse{UID: p.UID}
}
