package main

import (
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

const HELLO = "HELLO"
const PING = "PING"
const PONG = "PONG"
const ONLINE = "ONLINE"
const QUEUE = "QUEUE"
const UNQUEUE = "UNQUEUE"
const GAME = "GAME"
const PLAYERACTION = "PLAYERACTION"
const WORLD = "WORLD"
const EXITGAME = "EXITGAME"

type ClientMessage interface{}

type ServerMessage interface {
	Stringify() []byte
}

func ParseMessage(body []byte) ClientMessage {
	str := string(body)
	parts := strings.Split(str, ":")
	if len(parts) == 0 {
		return nil
	}

	messageType := parts[0]
	switch messageType {
	case HELLO:
		if parts[1] == "" {
			return NewHelloMessage(uuid.Nil)
		}

		uuid, err := uuid.Parse(parts[1])
		if nil != err {
			log.Println(err)
			return nil
		}

		return NewHelloMessage(uuid)
	case PONG:
		timestamp, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Println(err)
			return nil
		}

		return NewPongMessage(timestamp)
	case QUEUE:
		return NewQueueMessage()
	case UNQUEUE:
		return NewUnQueueMessage()
	case PLAYERACTION:
		x, err := strconv.Atoi(parts[1])
		if nil != err {
			log.Println(err)
			return nil
		}

		y, err := strconv.Atoi(parts[2])
		if nil != err {
			log.Println(err)
			return nil
		}

		return NewPlayerActionMessage(x, y)
	case EXITGAME:
		return NewExitGameMessage(parts[1])
	default:
		return nil
	}
}

type HelloMessage struct {
	playerId uuid.UUID
}

func NewHelloMessage(playerId uuid.UUID) HelloMessage {
	return HelloMessage{playerId: playerId}
}

func (message HelloMessage) Stringify() []byte {
	return []byte(HELLO + ":" + message.playerId.String())
}

type PingMessage struct {
	timestamp int64
	latency   int
}

func NewPingMessage(latency int) PingMessage {
	return PingMessage{
		timestamp: time.Now().UnixMilli(),
		latency:   latency,
	}
}

func (message PingMessage) Stringify() []byte {
	return []byte(PING + ":" + strconv.FormatInt(message.timestamp, 10) + ":" + strconv.Itoa(message.latency))
}

type PongMessage struct {
	timestamp int64
}

func NewPongMessage(timestamp int64) PongMessage {
	return PongMessage{
		timestamp: timestamp,
	}
}

type OnlineMessage struct {
	count int
}

func NewOnlineMessage(count int) OnlineMessage {
	return OnlineMessage{count: count}
}

func (message OnlineMessage) Stringify() []byte {
	return []byte(ONLINE + ":" + strconv.Itoa(message.count))
}

type QueueMessage struct{}

func NewQueueMessage() QueueMessage {
	return QueueMessage{}
}

type UnQueueMessage struct{}

func NewUnQueueMessage() UnQueueMessage {
	return UnQueueMessage{}
}

type PlayerActionMessage struct {
	x int
	y int
}

func NewPlayerActionMessage(x int, y int) PlayerActionMessage {
	return PlayerActionMessage{x: x, y: y}
}

type GameMessage struct {
	roomId uuid.UUID
}

func NewGameMessage(roomId uuid.UUID) GameMessage {
	return GameMessage{
		roomId: roomId,
	}
}

func (message GameMessage) Stringify() []byte {
	return []byte(GAME + ":" + message.roomId.String())
}

type ExitGameMessage struct {
	reason string
}

func NewExitGameMessage(reason string) ExitGameMessage {
	r, size := utf8.DecodeRuneInString(reason)
	if r == utf8.RuneError {
		log.Panicln(r)
	}

	s := string(unicode.ToUpper(r)) + reason[size:]

	return ExitGameMessage{
		reason: s,
	}
}

func (message ExitGameMessage) Stringify() []byte {
	return []byte(EXITGAME + ":" + message.reason)
}

type WorldMessage struct {
	posAX    int
	posAY    int
	posBX    int
	posBY    int
	posPuckX int
	posPuckY int
	countA   uint
	countB   uint
}

func NewWorldMessage(posA, posB, puck Position, countA uint, countB uint) WorldMessage {
	return WorldMessage{
		posAX:    posA.x,
		posAY:    posA.y,
		posBX:    posB.x,
		posBY:    posB.y,
		posPuckX: puck.x,
		posPuckY: puck.y,
		countA:   countA,
		countB:   countB,
	}
}

func (message WorldMessage) Stringify() []byte {
	return []byte(WORLD + ":" +
		strconv.Itoa(message.posAX) + ":" + strconv.Itoa(message.posAY) + ":" +
		strconv.Itoa(message.posBX) + ":" + strconv.Itoa(message.posBY) + ":" +
		strconv.Itoa(message.posPuckX) + ":" + strconv.Itoa(message.posPuckY) + ":" +
		strconv.Itoa(int(message.countA)) + ":" + strconv.Itoa(int(message.countB)))
}
