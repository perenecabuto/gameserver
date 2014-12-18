package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Broadcast chan []byte

type WebSocketManager interface {
	OnOpen(*websocket.Conn) []byte
	OnClose(*websocket.Conn)
	OnMessage(*websocket.Conn, []byte)
}

type WebSocketServer struct {
	broadcast Broadcast
	WebSocketManager
}

var (
	pingPeriod = time.Second * 1
	upgrader   = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

func NewWebSocketServer(manager WebSocketManager) *WebSocketServer {
	return &WebSocketServer{make(Broadcast), manager}
}

func (w *WebSocketServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(rw, r, nil)

	if err != nil {
		log.Println("WS connection failed: ", err)
		return
	}

	defer w.closeConnection(ws)

	go w.listenWSMessage(ws)
	go w.broadcastMessage(ws)

	w.broadcast <- w.OnOpen(ws)
	w.ping(ws)
}

func (w *WebSocketServer) listenWSMessage(ws *websocket.Conn) {
	//ws.SetReadLimit(maxMessageSize)
	//ws.SetReadDeadline(time.Now().Add(pongWait))
	//ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		log.Println("message_data:", message)
		w.broadcast <- message
	}
}

func (w *WebSocketServer) broadcastMessage(ws *websocket.Conn) {
	for {
		select {
		case message, ok := <-w.broadcast:
			if !ok {
				ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			log.Println("OnMessage:", message)
			w.OnMessage(ws, message)
		}
	}
}

func (w *WebSocketServer) closeConnection(ws *websocket.Conn) {
	defer ws.Close()
	w.OnClose(ws)
	log.Println("WS connection finished")
}

func (w *WebSocketServer) ping(ws *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//ws.SetWriteDeadline(time.Now().Add(writeWait))
			//return c.ws.WriteMessage(mt, payload)
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("WS error: ", err)
				return
			}
		}
	}
}
