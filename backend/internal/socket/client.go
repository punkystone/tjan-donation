package socket

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	writeWait = 10 * time.Second
	pongWait  = 60 * time.Second
	//nolint:mnd // allow
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 2048
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Error().Err(err)
				}
				return
			}
			writeMessage(message, c.Conn)
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func writeMessage(message []byte, connection *websocket.Conn) {
	log.Info().Msgf("recv: %s", string(message))
	w, err := connection.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	_, err = w.Write(message)
	if err != nil {
		log.Error().Err(err)
		return
	}
	if err := w.Close(); err != nil {
		return
	}
}
