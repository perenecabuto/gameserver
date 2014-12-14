package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/perenecabuto/gameserver/protobuf"
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
	if err != nil {
		log.Println("WS connection failed: ", err)
		return
	}
	defer closeConnection(ws)
	connections[ws] = true

	go ListenWSMessage(ws)
	go ListentAndBroadcastMessage(ws)

	data, _ := proto.Marshal(&protobuf.Chat{
		Name: proto.String(ws.RemoteAddr().String()),
		Text: proto.String("connected"),
	})
	broadcast <- data

	ping(ws)
}

func closeConnection(ws *websocket.Conn) {
	ws.Close()
	delete(connections, ws)

	data, _ := proto.Marshal(&protobuf.Chat{
		Name: proto.String(ws.RemoteAddr().String()),
		Text: proto.String("disconnected"),
	})

	broadcast <- data

	log.Println("WS connection finished")
}

func ListenWSMessage(ws *websocket.Conn) {
	//ws.SetReadLimit(maxMessageSize)
	//ws.SetReadDeadline(time.Now().Add(pongWait))
	//ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := ws.ReadMessage()
		println(message)
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

func ListentAndBroadcastMessage(ws *websocket.Conn) {
	for {
		select {
		case message, ok := <-broadcast:
			if !ok {
				ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			for ws := range connections {
				go ws.WriteMessage(websocket.BinaryMessage, message)
			}
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
