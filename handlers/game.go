package handlers

import (
	"errors"
	"log"

	"code.google.com/p/go-uuid/uuid"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"

	"github.com/perenecabuto/gameserver/game"
	"github.com/perenecabuto/gameserver/protobuf"
)

var (
	roomManager         = game.NewRoomManager()
	PlayerNotFoundError = errors.New("player not found")
)

type PlayerConnection struct {
	UUID uuid.UUID
	Conn *websocket.Conn
	X    *int32
	Y    *int32
}

func NewPlayerConnection(ws *websocket.Conn) *PlayerConnection {
	x, y := proto.Int32(100), proto.Int32(100)
	return &PlayerConnection{uuid.NewUUID(), ws, x, y}
}

func (p PlayerConnection) Id() string {
	return p.UUID.String()
}

func (p PlayerConnection) StatusMessage(action *protobuf.GameMessage_Action) *protobuf.GameMessage {
	return &protobuf.GameMessage{
		Id:       proto.String(p.Id()),
		Position: &protobuf.GameMessage_Position{X: p.X, Y: p.Y},
		Action:   action,
	}
}

type GameManager struct {
	room *game.Room
}

func GameHandler() *WebSocketServer {
	manager := &GameManager{roomManager.GetARoom()}
	return NewWebSocketServer(manager)
}

func (m *GameManager) OnOpen(ws *websocket.Conn) {
	player := NewPlayerConnection(ws)
	m.room = roomManager.GetARoom()
	m.room.Add(player)
	log.Println("Room", m.room)

	m.notifyClientNewPlayer(player)
	m.spawnRemotePlayers(ws)

	m.notifyAllPlayersAboutNewPlayer(player)
}

func (m *GameManager) OnMessage(ws *websocket.Conn, data []byte) {
	message := &protobuf.GameMessage{}
	proto.Unmarshal(data, message)
	player := m.room.GetById(*message.Id)
	if player == nil {
		return
	}

	log.Println(message)
	m.broadcastMessage(message)
}

func (m *GameManager) OnClose(ws *websocket.Conn) {
	player := m.playerByWS(ws)
	if player == nil {
		log.Println("Closing connection, now player found by WS!")
		return
	}

	log.Println("Removing player: ", player)
	m.room.Remove(player)

	message := player.StatusMessage(protobuf.GameMessage_DIE.Enum())
	m.broadcastMessage(message)
}

func (m *GameManager) playerByWS(ws *websocket.Conn) *PlayerConnection {
	for _, p := range m.room.Players {
		player := p.(*PlayerConnection)
		if player.Conn == ws {
			return player
		}
	}

	return nil
}

func (m *GameManager) broadcastMessage(gameMessage *protobuf.GameMessage) {
	message, _ := proto.Marshal(gameMessage)
	for _, p := range m.room.Players {
		player := p.(*PlayerConnection)
		go player.Conn.WriteMessage(websocket.BinaryMessage, message)
	}
}

func (m *GameManager) notifyClientNewPlayer(player *PlayerConnection) {
	message := player.StatusMessage(protobuf.GameMessage_CREATE.Enum())
	data, _ := proto.Marshal(message)
	player.Conn.WriteMessage(websocket.BinaryMessage, data)
}

func (m *GameManager) notifyAllPlayersAboutNewPlayer(player *PlayerConnection) {
	message := player.StatusMessage(protobuf.GameMessage_SPAWN.Enum())
	m.broadcastMessage(message)
}

// Notify new player about current players (with SPAWN status)
func (m *GameManager) spawnRemotePlayers(ws *websocket.Conn) {
	for _, p := range m.room.Players {
		player := p.(*PlayerConnection)
		if player.Conn == ws {
			continue
		}
		message := player.StatusMessage(protobuf.GameMessage_SPAWN.Enum())
		data, _ := proto.Marshal(message)
		ws.WriteMessage(websocket.BinaryMessage, data)
	}
}
