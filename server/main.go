package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const MESSAGE_TIMEOUT time.Duration = 10 * time.Second
const PING_RATE time.Duration = 1 * time.Second
const SERVER_CRT string = "server.crt"
const SERVER_KEY string = "server.key"

var pool = NewPool()

func main() {
	go pool.Matchmaking()

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		handler(writer, request)
	})

	_, errCrt := os.Stat(SERVER_CRT)
	_, errKey := os.Stat(SERVER_KEY)

	if errCrt == nil && errKey == nil {
		log.Println("Start HTTPS")
		log.Fatal(http.ListenAndServeTLS(":3001", SERVER_CRT, SERVER_KEY, nil))
	} else {
		log.Println("Start HTTP")
		log.Fatal(http.ListenAndServe(":3001", nil))
	}
}

func handler(writer http.ResponseWriter, request *http.Request) {
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return
	}

	go handleConnection(connection)
}

func handleConnection(connection *websocket.Conn) {
	conn := NewConnection(connection)

	player, success := conn.Handshake()
	if !success {
		connection.Close()

		return
	}

	read := make(chan ClientMessage)
	write := make(chan ServerMessage)

	player.SetExchange(read, write)

	go ping(player)
	go playerRead(conn, player)
	go playerWrite(conn, player)

	pool.UpdateOnline()
}

func handleDisconnection(conn *Connection) {
	conn.conn.Close()

	playerToDisconnect := pool.RemoveConnection(conn.id)
	if playerToDisconnect != nil {
		log.Println("Player disconnected: " + playerToDisconnect.id.String())
		pool.RemovePlayer(playerToDisconnect.id)
		pool.UpdateOnline()
	}
}

func ping(player *Player) {
	lastPing := time.Now()

	for {
		if time.Since(lastPing) < PING_RATE {
			continue
		}

		player.GetWrite() <- NewPingMessage(player.latency)

		lastPing = time.Now()
	}
}

func playerRead(conn *Connection, player *Player) {
	for {
		message, err := conn.ReadMessage()

		if err != nil {
			handleDisconnection(conn)

			return
		}

		if message == nil {
			continue
		}

		switch val := message.(type) {
		case QueueMessage:
			pool.QueuePlayer(player)
		case UnQueueMessage:
			pool.UnQueuePlayer(player)
		case PongMessage:
			player.CalcLatency(val)
		case ExitGameMessage:
			pool.DeleteRoom(player.currentRoomId, errors.New(val.reason))
		case PlayerActionMessage:
			player.GetInput() <- message
		}
	}
}

func playerWrite(conn *Connection, player *Player) {
	for {
		for message := range player.GetWrite() {
			if message == nil {
				continue
			}

			if err := conn.WriteMessage(message); err != nil {
				handleDisconnection(conn)

				return
			}
		}
	}
}
