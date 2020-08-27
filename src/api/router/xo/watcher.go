package xo

import (
	"api/app"
	"api/bus"
)

type Watcher struct {
	*app.Participant
}

func (g *Game) NewWatcher(participant *app.Participant) (*Watcher, error) {
	watcher := &Watcher{
		Participant: participant,
	}
	participant.Meta = "watcher"
	participant.Bus.NewCallback(app.ConnectEventType, watcher.onUserConnect, g.Board)
	g.Bus.NewCallback(BoardChangesEventType, watcher.onBoardChanged, nil)
	g.Bus.NewCallback(EndGameEventType, watcher.onEndGame, nil)
	return watcher, nil
}

func (watcher *Watcher) onBoardChanged(args bus.CallbackArgs) error {
	watcher.Absorber <- BuildBoardStateMessage(args.Initiator.(*Board))
	return nil
}

func (watcher *Watcher) onEndGame(args bus.CallbackArgs) error {
	watcher.Absorber <- BuildWatcherEndGameMessage(args.Event.GetPayload().(SignType))
	return nil
}

func (watcher *Watcher) onUserConnect(args bus.CallbackArgs) error {
	watcher.Absorber <- BuildBoardStateMessage(args.Meta.(*Board))
	return nil
}
