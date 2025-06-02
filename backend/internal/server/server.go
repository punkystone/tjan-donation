package server

import (
	"log"
	"net/http"
	"time"
	"tjan-donation/internal/socket"

	"github.com/gorilla/websocket"
)

const timeout = 5
const ReadBufferSize = 1024
const WriteBufferSize = 1024
const sendChannelSize = 256

type Handler struct {
	hub *socket.Hub
}

func (handler *Handler) handle(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  ReadBufferSize,
		WriteBufferSize: WriteBufferSize,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &socket.Client{Hub: handler.hub, Conn: conn, Send: make(chan []byte, sendChannelSize)}
	client.Hub.Register <- client

	go client.Write()
}

func StartServer(hub *socket.Hub) error {
	handler := &Handler{
		hub: hub,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.handle)
	server := &http.Server{
		Addr:              ":80",
		Handler:           mux,
		ReadHeaderTimeout: timeout * time.Second,
	}

	return server.ListenAndServe()
}
