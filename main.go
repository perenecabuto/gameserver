package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/perenecabuto/gameserver/protobuf"
)

type Configuration struct {
	Address         string `json:"addr"`
	Port            string `json:"port"`
	CertificateFile string `json:"certFile"`
	KeyFile         string `json:"keyFile"`
	UseTLS          bool   `json:"useTLS"`
}

type Broadcast chan []byte
type ConnectionMap map[*websocket.Conn]bool

var (
	configFile = flag.String("config", "config.json", "configuration file")

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	pingPeriod = time.Second * 1
)

type WebSocketServer struct {
	connections ConnectionMap
	broadcast   Broadcast
	OnOpen      func(*websocket.Conn) []byte
}

func main() {
	flag.Parse()

	config, err := ReadConfiguration()
	if err != nil {
		log.Fatal("error starting: ", err)
	}

	serverAddress := fmt.Sprintf("%s:%s", config.Address, config.Port)

	Routes()

	if config.UseTLS {
		log.Println(fmt.Sprintf("Starting Secure WebSocket server at https://%s", serverAddress))
		log.Fatal(http.ListenAndServeTLS(serverAddress, config.CertificateFile, config.KeyFile, nil))
	} else {
		log.Println(fmt.Sprintf("Starting WebSocket server at http://%s", serverAddress))
		log.Fatal(http.ListenAndServe(serverAddress, nil))
	}
}

func ReadConfiguration() (conf *Configuration, err error) {
	conf = &Configuration{}

	var file *os.File
	if file, err = os.Open(*configFile); err == nil {
		log.Println("Reading config:", *configFile)

		decoder := json.NewDecoder(file)

		if err = decoder.Decode(&conf); err != nil {
			return
		}
	}

	return
}

func Routes() {
	router := mux.NewRouter()

	router.PathPrefix("/protobuf/").Handler(http.StripPrefix("/protobuf/", http.FileServer(http.Dir("protobuf/"))))

	router.Handle("/ws/chat", &WebSocketServer{make(ConnectionMap), make(Broadcast), ChatOnOpen})
	router.Handle("/ws/game", &WebSocketServer{make(ConnectionMap), make(Broadcast), GameOnOpen})
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("webroot/")))

	http.Handle("/", router)
}

func (w *WebSocketServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(rw, r, nil)

	if err != nil {
		log.Println("WS connection failed: ", err)
		return
	}

	defer w.closeConnection(ws)
	w.connections[ws] = true

	go w.ListenWSMessage(ws)
	go w.BroadcastMessage(ws)

	w.broadcast <- w.OnOpen(ws)
	w.ping(ws)
}

func ChatOnOpen(ws *websocket.Conn) []byte {
	data, _ := proto.Marshal(&protobuf.ChatMessage{
		Name:        proto.String(ws.RemoteAddr().String()),
		Text:        proto.String(fmt.Sprintf("[%p] connected", ws)),
		MessageType: protobuf.ChatMessage_CONNECTION.Enum(),
	})

	return data
}

func GameOnOpen(ws *websocket.Conn) []byte {
	data, _ := proto.Marshal(&protobuf.GameMessage{
		Id:       proto.Int32(10),
		Position: &protobuf.GameMessage_Position{X: proto.Int32(1), Y: proto.Int32(10)},
		Action:   protobuf.GameMessage_SPAWN.Enum(),
	})

	return data
}

func (w *WebSocketServer) closeConnection(ws *websocket.Conn) {
	ws.Close()
	delete(w.connections, ws)

	data, _ := proto.Marshal(&protobuf.ChatMessage{
		Name:        proto.String(ws.RemoteAddr().String()),
		Text:        proto.String(fmt.Sprintf("[%p] disconnected", ws)),
		MessageType: protobuf.ChatMessage_DISCONNECTION.Enum(),
	})

	w.broadcast <- data

	log.Println("WS connection finished")
}

func (w *WebSocketServer) ListenWSMessage(ws *websocket.Conn) {
	//ws.SetReadLimit(maxMessageSize)
	//ws.SetReadDeadline(time.Now().Add(pongWait))
	//ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := ws.ReadMessage()

		println(message)

		if err != nil {
			break
		}

		log.Println("message_data:", message)

		w.broadcast <- message
	}
}

func (w *WebSocketServer) BroadcastMessage(ws *websocket.Conn) {
	for {
		select {
		case message, ok := <-w.broadcast:
			if !ok {
				ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			for ws := range w.connections {
				go ws.WriteMessage(websocket.BinaryMessage, message)
			}
		}
	}
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
