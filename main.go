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
}

var (
	configFile = flag.String("config", "config.json", "configuration file")

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

	config, err := ReadConfiguration()

	if err != nil {
		log.Fatal("error starting: ", err)
	}

	serverAddress := fmt.Sprintf("%s:%s", config.Address, config.Port)

	Routes()

	if config.CertificateFile != "" && config.KeyFile != "" {
		log.Println(fmt.Sprintf("Starting Secure WebSocket server at https://%s", serverAddress))
		log.Fatal(http.ListenAndServeTLS(serverAddress, config.CertificateFile, config.KeyFile, nil))
	} else {
		log.Println(fmt.Sprintf("Starting WebSocket server at http://%s", serverAddress))
		log.Fatal(http.ListenAndServe(serverAddress, nil))
	}
}

func ReadConfiguration() (*Configuration, error) {
	file, err := os.Open(*configFile)
	configuration := Configuration{}

	if err == nil {
		log.Println("Reading config:", *configFile)

		decoder := json.NewDecoder(file)
		err := decoder.Decode(&configuration)

		if err != nil {
			log.Fatal("json error: ", err)
		}
	}

	return &configuration, err
}

func Routes() {
	router := mux.NewRouter()

	router.PathPrefix("/protobuf/").Handler(http.StripPrefix("/protobuf/", http.FileServer(http.Dir("protobuf/"))))
	router.HandleFunc("/ws", WebSocketServer)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("webroot/")))

	http.Handle("/", router)
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
	go BroadcastMessage(ws)

	data, _ := proto.Marshal(&protobuf.ChatMessage{
		Name:        proto.String(ws.RemoteAddr().String()),
		Text:        proto.String(fmt.Sprintf("[%p] connected", ws)),
		MessageType: protobuf.ChatMessage_CONNECTION.Enum(),
	})

	broadcast <- data

	ping(ws)
}

func closeConnection(ws *websocket.Conn) {
	ws.Close()
	delete(connections, ws)

	data, _ := proto.Marshal(&protobuf.ChatMessage{
		Name:        proto.String(ws.RemoteAddr().String()),
		Text:        proto.String(fmt.Sprintf("[%p] disconnected", ws)),
		MessageType: protobuf.ChatMessage_DISCONNECTION.Enum(),
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

		receivedChat := new(protobuf.ChatMessage)

		err = proto.Unmarshal(message, receivedChat)

		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}

		log.Println("[encoded_message] data:", message)
		log.Println("[decoded_message] user: ", receivedChat.GetName(), " text: ", receivedChat.GetText())

		broadcast <- message
	}
}

func BroadcastMessage(ws *websocket.Conn) {
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
