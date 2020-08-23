package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Transport struct {
	Reader  chan []byte
	Writer  chan []byte
	Context context.Context
	Cancel  context.CancelFunc
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = pongWait / 2
	maxMessageSize = 512
)

func NewWsTransport(e echo.Context) (*Transport, error) {
	reader := make(chan []byte)
	writer := make(chan []byte)

	conn, err := Upgrader.Upgrade(e.Response(), e.Request(), nil)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}
	e.Logger().Debug("Initialize client complete")
	ctx, cancel := context.WithCancel(context.Background())
	transport := &Transport{Reader: reader, Writer: writer, Context: ctx, Cancel: cancel}

	go func() { // reader binding
		var err error

		defer func() {
			if err != nil {
				e.Logger().Error(err)
			}
			transport.Cancel()
			conn.Close()
		}()

		conn.SetReadLimit(maxMessageSize)
		if err = conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return
		}
		conn.SetPongHandler(func(string) error { err = conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
		for {
			select {
			case <-transport.Context.Done():
				return
			default:
				_, message, err := conn.ReadMessage()
				if err != nil {
					return
				}
				reader <- message // TODO: deserialize
			}
		}
	}()

	go func() { // writer binding
		var err error

		defer func() {
			if err != nil {
				e.Logger().Error(err)
			}
			transport.Cancel()
			conn.Close()
		}()

		ticker := time.NewTicker(pingPeriod)
		for {
			select {
			case message := <-transport.Writer:
				if err = conn.WriteMessage(websocket.TextMessage, message); err != nil {
					return
				} // TODO: serialize
			case <-ticker.C:
				if err = conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
					return
				}
				if err = conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			case <-transport.Context.Done():
				return
			}
		}
	}()

	return transport, nil
}
