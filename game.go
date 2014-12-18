package main

import (
	"log"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/perenecabuto/gameserver/protobuf"
)

type PlayerConnection struct {
	Conn        *websocket.Conn
	LastMessage *protobuf.GameMessage
}

type GameManager struct {
	currentId int32
	players   map[int32]*PlayerConnection
}

func NewGameServer() *WebSocketServer {
	manager := &GameManager{0, make(map[int32]*PlayerConnection)}
	return NewWebSocketServer(manager)
}

func (m *GameManager) OnOpen(ws *websocket.Conn) []byte {
	defer func() { m.currentId++ }()

	message := &protobuf.GameMessage{
		Id:       proto.Int32(m.currentId),
		Position: &protobuf.GameMessage_Position{X: proto.Int32(1), Y: proto.Int32(10)},
		Action:   protobuf.GameMessage_SPAWN.Enum(),
	}
	data, _ := proto.Marshal(message)
	message.Action = protobuf.GameMessage_IDLE.Enum()
	m.players[*message.Id] = &PlayerConnection{ws, message}

	return data
}

func (m *GameManager) OnMessage(ws *websocket.Conn, message []byte) {
	data := &protobuf.GameMessage{}
	proto.Unmarshal(message, data)
	player, ok := m.players[*data.Id]
	if !ok {
		return
	}

	player.LastMessage = data
	newMessage, _ := proto.Marshal(data)
	m.broadcastMessage(newMessage)
}

func (m *GameManager) OnClose(ws *websocket.Conn) {
	player := m.playerByWS(ws)
	player.LastMessage.Action = protobuf.GameMessage_DEAD.Enum()
	message, _ := proto.Marshal(player.LastMessage)

	log.Println("Removing ws: ", player.LastMessage)
	delete(m.players, *player.LastMessage.Id)

	m.broadcastMessage(message)
}

func (m *GameManager) playerByWS(ws *websocket.Conn) *PlayerConnection {
	for _, player := range m.players {
		if player.Conn == ws {
			return player
		}
	}

	return nil
}

func (m *GameManager) broadcastMessage(message []byte) {
	for _, player := range m.players {
		player.Conn.WriteMessage(websocket.BinaryMessage, message)
	}
}
