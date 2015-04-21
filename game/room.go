package game

import "log"

const (
	MaxPlayersPerRoom = 4
)

type Player interface {
	Id() string
}

type Room struct {
	Players map[string]Player
}

func NewRoom() *Room {
	log.Println("Create New Room")
	players := make(map[string]Player)
	return &Room{players}
}

func (r *Room) Add(p Player) {
	if r.Crowded() {
		panic("MaxPlayersPerRoom reached")
	}
	r.Players[p.Id()] = p
}

func (r *Room) Remove(p Player) {
	if p != nil && r.Players[p.Id()] != nil {
		delete(r.Players, p.Id())
	}
}

func (r *Room) Crowded() bool {
	log.Println(&r, len(r.Players))
	return len(r.Players) >= MaxPlayersPerRoom
}

func (r *Room) GetById(id string) Player {
	return r.Players[id]
}

type RoomManager struct {
	Rooms []*Room
}

func NewRoomManager() *RoomManager {
	rooms := []*Room{NewRoom()}
	return &RoomManager{rooms}
}

func (m *RoomManager) GetARoom() *Room {
	room := m.Rooms[len(m.Rooms)-1]

	if room.Crowded() {
		room = NewRoom()
		m.Rooms[len(m.Rooms)] = room
	}

	return room
}
