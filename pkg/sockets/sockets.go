package sockets

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CloseWebSocket(ws *websocket.Conn, code int, reason string) {
	err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, reason))

	if err != nil {
		log.Println(err)
	}

	ws.Close()
}
