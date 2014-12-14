// Code generated by protoc-gen-go.
// source: chat.proto
// DO NOT EDIT!

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	chat.proto
	object.proto

It has these top-level messages:
	ChatMessage
*/
package protobuf

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type ChatMessage_MessageType int32

const (
	ChatMessage_MESSAGE       ChatMessage_MessageType = 0
	ChatMessage_CONNECTION    ChatMessage_MessageType = 1
	ChatMessage_DISCONNECTION ChatMessage_MessageType = 2
)

var ChatMessage_MessageType_name = map[int32]string{
	0: "MESSAGE",
	1: "CONNECTION",
	2: "DISCONNECTION",
}
var ChatMessage_MessageType_value = map[string]int32{
	"MESSAGE":       0,
	"CONNECTION":    1,
	"DISCONNECTION": 2,
}

func (x ChatMessage_MessageType) Enum() *ChatMessage_MessageType {
	p := new(ChatMessage_MessageType)
	*p = x
	return p
}
func (x ChatMessage_MessageType) String() string {
	return proto.EnumName(ChatMessage_MessageType_name, int32(x))
}
func (x *ChatMessage_MessageType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ChatMessage_MessageType_value, data, "ChatMessage_MessageType")
	if err != nil {
		return err
	}
	*x = ChatMessage_MessageType(value)
	return nil
}

type ChatMessage struct {
	Name             *string                  `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	Text             *string                  `protobuf:"bytes,2,req,name=text" json:"text,omitempty"`
	MessageType      *ChatMessage_MessageType `protobuf:"varint,3,req,name=messageType,enum=protobuf.ChatMessage_MessageType,def=0" json:"messageType,omitempty"`
	XXX_unrecognized []byte                   `json:"-"`
}

func (m *ChatMessage) Reset()         { *m = ChatMessage{} }
func (m *ChatMessage) String() string { return proto.CompactTextString(m) }
func (*ChatMessage) ProtoMessage()    {}

const Default_ChatMessage_MessageType ChatMessage_MessageType = ChatMessage_MESSAGE

func (m *ChatMessage) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *ChatMessage) GetText() string {
	if m != nil && m.Text != nil {
		return *m.Text
	}
	return ""
}

func (m *ChatMessage) GetMessageType() ChatMessage_MessageType {
	if m != nil && m.MessageType != nil {
		return *m.MessageType
	}
	return Default_ChatMessage_MessageType
}

func init() {
	proto.RegisterEnum("protobuf.ChatMessage_MessageType", ChatMessage_MessageType_name, ChatMessage_MessageType_value)
}
