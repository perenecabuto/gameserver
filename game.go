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

func (m *GameManager) OnOpen(ws *websocket.Conn) {
	defer func() { m.currentId++ }()

	// Send new player his start data
	message := &protobuf.GameMessage{
		Id:       proto.Int32(m.currentId),
		Position: &protobuf.GameMessage_Position{X: proto.Int32(100), Y: proto.Int32(200)},
		Action:   protobuf.GameMessage_CREATE.Enum(),
	}

	data, _ := proto.Marshal(message)
	ws.WriteMessage(websocket.BinaryMessage, data)

	m.spawnRemotePlayers(ws)

	// Add new player to players collection
	m.players[*message.Id] = &PlayerConnection{ws, message}

	// Notify all players about new player
	message = proto.Clone(message).(*protobuf.GameMessage)
	message.Action = protobuf.GameMessage_SPAWN.Enum()
	m.broadcastMessage(message)
}

func (m *GameManager) OnMessage(ws *websocket.Conn, data []byte) {
	message := &protobuf.GameMessage{}
	proto.Unmarshal(data, message)
	player, ok := m.players[*message.Id]
	if !ok {
		return
	}

	player.LastMessage = message
	m.broadcastMessage(message)
}

func (m *GameManager) OnClose(ws *websocket.Conn) {
	player := m.playerByWS(ws)
	player.LastMessage.Action = protobuf.GameMessage_DIE.Enum()

	log.Println("Removing ws: ", player.LastMessage)
	delete(m.players, *player.LastMessage.Id)

	m.broadcastMessage(player.LastMessage)
}

func (m *GameManager) playerByWS(ws *websocket.Conn) *PlayerConnection {
	for _, player := range m.players {
		if player.Conn == ws {
			return player
		}
	}

	return nil
}

func (m *GameManager) broadcastMessage(gameMessage *protobuf.GameMessage) {
	message, _ := proto.Marshal(gameMessage)
	for _, player := range m.players {
		go player.Conn.WriteMessage(websocket.BinaryMessage, message)
	}
}

// Notify new player about current players (with SPAWN status)
func (m *GameManager) spawnRemotePlayers(ws *websocket.Conn) {
	for _, player := range m.players {
		if player.Conn == ws {
			continue
		}
		m := proto.Clone(player.LastMessage).(*protobuf.GameMessage)
		m.Action = protobuf.GameMessage_SPAWN.Enum()
		data, _ := proto.Marshal(m)
		ws.WriteMessage(websocket.BinaryMessage, data)
	}
}
