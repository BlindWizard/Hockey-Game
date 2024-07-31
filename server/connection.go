package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Connection struct {
	id   uuid.UUID
	conn *websocket.Conn
}

func NewConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		id:   uuid.New(),
		conn: conn,
	}
}

func (connection *Connection) ReadMessage() (ClientMessage, error) {
	connection.conn.SetReadDeadline(time.Now().Add(MESSAGE_TIMEOUT))
	_, bytes, err := connection.conn.ReadMessage()
	if err != nil {
		log.Println(err)

		return nil, err
	}

	message := ParseMessage(bytes)
	if nil == message {
		return nil, nil
	}

	return message, nil
}

func (connection *Connection) WriteMessage(message ServerMessage) error {
	connection.conn.SetWriteDeadline(time.Now().Add(MESSAGE_TIMEOUT))
	err := connection.conn.WriteMessage(websocket.TextMessage, message.Stringify())
	if err != nil {
		log.Println(err)

		return err
	}

	return nil
}

func (conn *Connection) Handshake() (*Player, bool) {
	handshakeStarted := time.Now()

	var clientHello HelloMessage

	for {
		message, err := conn.ReadMessage()
		if err != nil {
			return nil, false
		}

		if message == nil {
			return nil, false
		}

		value, ok := message.(HelloMessage)
		if !ok {
			continue
		}

		if time.Since(handshakeStarted) > MESSAGE_TIMEOUT {
			return nil, false
		}

		clientHello = value
		break
	}

	playerId := clientHello.playerId
	player := pool.NewPlayer(playerId, conn)

	log.Println("New player connected " + player.id.String())

	if err := conn.WriteMessage(NewHelloMessage(player.id)); err != nil {
		pool.RemovePlayer(playerId)

		return nil, false
	}

	return player, true
}
