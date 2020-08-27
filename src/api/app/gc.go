package app

import "time"

const roomTTL = 10 * time.Minute
const participantTTL = 5 * time.Minute

func (app *Context) RunGC() {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-app.context.Done():
			return
		case <-t.C:
			for roomUID, roomIdx := range app.roomIndex {
				room := app.rooms[roomIdx]
				if room.lastModified.Add(roomTTL).Before(time.Now()) {
					app.DeleteRoom(roomUID)
				}
			}
		}
	}
}

func (room *RoomContext) RunGC() {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-room.context.Done():
			return
		case <-t.C:
			for UID, participantIdx := range room.participantIndex {
				participant := room.Participants[participantIdx]
				if participant.lastModified.Add(participantTTL).Before(time.Now()) {
					room.DeleteParticipant(UID)
				}
			}
		}
	}
}
