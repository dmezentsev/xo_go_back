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
			for roomUID, room := range app.rooms {
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
			for UID, p := range room.Participants {
				if p.lastModified.Add(participantTTL).Before(time.Now()) {
					room.DeleteParticipant(UID)
				}
			}
		}
	}
}
