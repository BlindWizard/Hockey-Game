package main

import (
	"errors"
	"log"
	"math/rand"
	"slices"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Pool struct {
	connections map[uuid.UUID]*Connection
	players     map[uuid.UUID]*Player
	queue       []*Player
	rooms       map[uuid.UUID]*Room
}

func NewPool() *Pool {
	return &Pool{
		make(map[uuid.UUID]*Connection),
		make(map[uuid.UUID]*Player),
		make([]*Player, 0, 100),
		make(map[uuid.UUID]*Room),
	}
}

func (pool *Pool) NewPlayer(playerId uuid.UUID, conn *Connection) *Player {
	if playerId == uuid.Nil {
		playerId = uuid.New()
	}

	player := NewPlayer(playerId, conn.id)
	pool.connections[conn.id] = conn
	pool.players[player.id] = player

	return player
}

func (pool *Pool) CreateRoom(playerA *Player, playerB *Player) *Room {
	var room *Room
	if rand.Intn(2) == 1 {
		room = NewRoom(
			playerA,
			playerB,
		)
	} else {
		room = NewRoom(
			playerB,
			playerA,
		)
	}

	playerA.currentRoomId = room.id
	playerB.currentRoomId = room.id

	pool.rooms[room.id] = room

	log.Println("Start game: " + room.id.String())

	return room
}

func (pool *Pool) RemoveConnection(connId uuid.UUID) *Player {
	delete(pool.connections, connId)
	for _, player := range pool.players {
		if connId == player.connectionId {
			return player
		}
	}

	return nil
}

func (pool *Pool) RemovePlayer(id uuid.UUID) {
	player, exists := pool.players[id]
	if !exists {
		return
	}

	delete(pool.players, id)
	delete(pool.connections, player.connectionId)

	pool.UnQueuePlayer(player)
	if player.currentRoomId != uuid.Nil {
		pool.DeleteRoom(player.currentRoomId, errors.New("player disconnected"))
	}
}

func (pool *Pool) DeleteRoom(roomId uuid.UUID, reason error) {
	room, exists := pool.rooms[roomId]
	if !exists {
		return
	}

	go room.Close(reason)

	for _, player := range []*Player{room.playerA, room.playerB} {
		player.GetWrite() <- NewExitGameMessage(reason.Error())
	}

	delete(pool.rooms, roomId)
}

func (pool *Pool) QueuePlayer(player *Player) {
	for _, queuedPlayer := range pool.queue {
		if queuedPlayer.id == player.id {
			return
		}
	}

	pool.queue = append(pool.queue, player)

	log.Println("Current queue: " + strconv.Itoa(len(pool.queue)))
}

func (pool *Pool) UnQueuePlayer(player *Player) {
	pool.queue = slices.DeleteFunc(pool.queue, func(other *Player) bool {
		return player.id == other.id
	})

	log.Println("Current queue: " + strconv.Itoa(len(pool.queue)))
}

func (pool *Pool) UpdateOnline() {
	log.Println("Current online: " + strconv.Itoa(len(pool.players)))

	for _, player := range pool.players {
		player.GetWrite() <- NewOnlineMessage(len(pool.players))
	}
}

func (pool *Pool) Matchmaking() {
	for {
		for {
			if cap(pool.queue) < 2 || len(pool.queue) < 2 {
				break
			}

			playersFromQueue := pool.queue[0:2]

			if playersFromQueue[0] == nil || playersFromQueue[1] == nil {
				break
			}

			pool.queue = pool.queue[2:]

			room := pool.CreateRoom(playersFromQueue[0], playersFromQueue[1])
			for _, player := range playersFromQueue {
				player.writeChan <- NewGameMessage(room.id)
			}

			log.Println("Current queue: " + strconv.Itoa(len(pool.queue)))
			log.Println("Current games: " + strconv.Itoa(len(pool.rooms)))

			go room.RunGame()
		}

		time.Sleep(time.Second)
	}
}
