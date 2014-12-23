package main

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/perenecabuto/gameserver/protobuf"
)

type ConnectionMap map[*websocket.Conn]bool
type ChatManager struct {
	connections ConnectionMap
}

func NewChatServer() *WebSocketServer {
	manager := &ChatManager{make(ConnectionMap)}
	return NewWebSocketServer(manager)
}

func (c *ChatManager) OnOpen(ws *websocket.Conn) {
	c.connections[ws] = true

	message := &protobuf.ChatMessage{
		Name:        proto.String(ws.RemoteAddr().String()),
		Text:        proto.String(fmt.Sprintf("[%p] connected", ws)),
		MessageType: protobuf.ChatMessage_CONNECTION.Enum(),
	}

	c.sendChatMessage(message)
}

func (c *ChatManager) OnMessage(ws *websocket.Conn, message []byte) {
	c.sendMessage(message)
}

func (c *ChatManager) OnClose(ws *websocket.Conn) {
	delete(c.connections, ws)

	message := &protobuf.ChatMessage{
		Name:        proto.String(ws.RemoteAddr().String()),
		Text:        proto.String(fmt.Sprintf("[%p] disconnected", ws)),
		MessageType: protobuf.ChatMessage_DISCONNECTION.Enum(),
	}

	c.sendChatMessage(message)
}

func (c *ChatManager) sendChatMessage(chatMessage *protobuf.ChatMessage) {
	message, _ := proto.Marshal(chatMessage)
	c.sendMessage(message)
}

func (c *ChatManager) sendMessage(message []byte) {
	for ws := range c.connections {
		go ws.WriteMessage(websocket.BinaryMessage, message)
	}
}
