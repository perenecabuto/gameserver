package main

import (
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/perenecabuto/gameserver/protobuf"
	"log"
	"net/http"
	"time"
)

var (
	addr = flag.String("addr", ":3000", "http service address")

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	pingPeriod  = time.Second * 1
	broadcast   = make(chan []byte)
	connections = make(map[*websocket.Conn]bool)
)

func main() {
	flag.Parse()
	log.Println("Staring WebSocket server at ", *addr)

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", WebSocketServer)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func WebSocketServer(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	connections[ws] = true

	if err != nil {
		log.Println("WS connection failed: ", err)
		return
	}

	defer func() {
		ws.Close()
		delete(connections, ws)
		broadcastMessage([]byte(fmt.Sprintf("[%p] %s disconnected", ws, ws.RemoteAddr())))
		log.Println("WS connection finished")
	}()

	broadcastMessage([]byte(fmt.Sprintf("[%p] %s connected", ws, ws.RemoteAddr())))

	go sendMessage(ws)
	go readMessage(ws)

	ping(ws)
}

func readMessage(ws *websocket.Conn) {
	//ws.SetReadLimit(maxMessageSize)
	//ws.SetReadDeadline(time.Now().Add(pongWait))
	//ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := ws.ReadMessage()

		if err != nil {
			break
		}

		receivedChat := new(protobuf.Chat)
		err = proto.Unmarshal(message, receivedChat)

		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}

		log.Println("[encoded_message] data:", message)
		log.Println("[decoded_message] user: ", receivedChat.GetName(), " text: ", receivedChat.GetText())

		broadcast <- message
	}
}

func sendMessage(ws *websocket.Conn) {
	for {
		select {
		case message, ok := <-broadcast:
			if !ok {
				ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			broadcastMessage(message)
		}
	}
}

func broadcastMessage(message []byte) {
	for ws := range connections {
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func ping(ws *websocket.Conn) {
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
