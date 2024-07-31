package main

import (
	"time"

	"github.com/google/uuid"
)

type Player struct {
	id            uuid.UUID
	connectionId  uuid.UUID
	currentRoomId uuid.UUID
	latency       int
	inputChan     chan ClientMessage
	writeChan     chan ServerMessage
}

func NewPlayer(id uuid.UUID, connectionId uuid.UUID) *Player {
	return &Player{
		id:           id,
		connectionId: connectionId,
		latency:      0,
	}
}

func (player *Player) SetCurrentRoomId(id uuid.UUID) {
	player.currentRoomId = id
}

func (player *Player) SetExchange(read chan ClientMessage, write chan ServerMessage) {
	player.inputChan = read
	player.writeChan = write
}

func (player *Player) CalcLatency(ping PongMessage) {
	prevTime := time.UnixMilli(ping.timestamp)
	player.latency = int(time.Since(prevTime).Milliseconds() / 2)
}

func (player *Player) GetInput() chan ClientMessage {
	return player.inputChan
}

func (player *Player) GetWrite() chan ServerMessage {
	return player.writeChan
}
