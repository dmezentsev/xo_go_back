package protocol

import "api/app"

type RoomResponse struct {
	UID app.UIDType `json:"uid"`
}

func RoomSerialize(r *app.RoomContext) RoomResponse {
	return RoomResponse{UID: r.UID}
}
