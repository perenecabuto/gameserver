package main

import (
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

	m.sendPlayerStatus(ws)

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
}

func (m *GameManager) OnClose(ws *websocket.Conn) {
	//delete(players, id)
}

func (m *GameManager) sendPlayerStatus(ws *websocket.Conn) {
	for _, player := range m.players {
		message, _ := proto.Marshal(player.LastMessage)
		ws.WriteMessage(websocket.BinaryMessage, message)
	}
}
