package main

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	id      uuid.UUID
	playerA *Player
	playerB *Player
	world   *World
	exit    chan error
}

func NewRoom(playerA *Player, playerB *Player) *Room {
	return &Room{
		uuid.New(),
		playerA,
		playerB,
		NewWorld(),
		make(chan error),
	}
}

func (room *Room) Close(reason error) {
	log.Println("Close game: " + room.id.String() + " - " + reason.Error())
	room.exit <- reason
}

func (room *Room) RunGame() {
	updateTicker := time.NewTicker(time.Duration(PHYSICS_CYCLE) * time.Millisecond)
	broadcastTicker := time.NewTicker(time.Duration(NETWORK_CYCLE) * time.Millisecond)

	timerA := time.Now()
	timerB := time.Now()
	timerWorld := time.Now()

	for {
		select {
		case msg := <-room.playerA.GetInput():
			if time.Since(timerA).Milliseconds() < PLAYER_MESSAGE_THROTTLE {
				continue
			}

			message, ok := msg.(PlayerActionMessage)
			if !ok {
				continue
			}

			room.handlePlayerA(message)
			timerA = time.Now()
		case msg := <-room.playerB.GetInput():
			if time.Since(timerB).Milliseconds() < PLAYER_MESSAGE_THROTTLE {
				continue
			}

			message, ok := msg.(PlayerActionMessage)
			if !ok {
				continue
			}

			room.handlePlayerB(message)
			timerB = time.Now()
		case <-updateTicker.C:
			room.updateWorldState(time.Since(timerWorld))
			timerWorld = time.Now()
		case <-broadcastTicker.C:
			room.broadcastWorldState()
		case <-room.exit:
			updateTicker.Stop()
			broadcastTicker.Stop()

			return
		}
	}
}

func (room *Room) handlePlayerA(message PlayerActionMessage) {
	newPosition := validatePosition(NewPosition(message.x, message.y))

	room.world.positionPrevA.x = room.world.positionA.x
	room.world.positionPrevA.y = room.world.positionA.y
	room.world.positionA.x = newPosition.x
	room.world.positionA.y = newPosition.y
}

func (room *Room) handlePlayerB(message PlayerActionMessage) {
	pos := validatePosition(NewPosition(message.x, message.y))
	flip := FlipPosition(pos)

	room.world.positionPrevB.x = room.world.positionB.x
	room.world.positionPrevB.y = room.world.positionB.y
	room.world.positionB.x = flip.x
	room.world.positionB.y = flip.y
}

func (room *Room) updateWorldState(deltaT time.Duration) {
	newPuckPosition := NewVector(
		room.world.positionPuck.x+room.world.magnitudePuck.x*float64(deltaT.Milliseconds()),
		room.world.positionPuck.y+room.world.magnitudePuck.y*float64(deltaT.Milliseconds()),
	)

	collision, hitPoint, wallNormal := detectWallHit(room.world.walls, newPuckPosition, room.world.positionPuck)
	if collision {
		room.world.magnitudePuck = hitWall(room.world.magnitudePuck, wallNormal)
		room.world.positionPuck = hitPoint

		return
	}

	ok, position, magnitude := detectPlayerHit(
		room.world.positionA,
		room.world.positionPrevA,
		newPuckPosition,
		room.world.positionPuck,
		room.world.magnitudePuck,
		deltaT,
	)

	if ok {
		room.world.positionPuck = position
		room.world.magnitudePuck = validateMagnitute(magnitude)

		return
	}

	ok, position, magnitude = detectPlayerHit(
		room.world.positionB,
		room.world.positionPrevB,
		newPuckPosition,
		room.world.positionPuck,
		room.world.magnitudePuck,
		deltaT,
	)

	if ok {
		room.world.positionPuck = position
		room.world.magnitudePuck = validateMagnitute(magnitude)

		return
	}

	room.world.positionPuck = newPuckPosition
	room.world.magnitudePuck = MultiplyVectorNumber(room.world.magnitudePuck, 1-DEFAULT_DRAG)

	goal, side := detectGoal(*room.world.positionPuck)
	if goal {
		room.world.magnitudePuck = NewVector(0, 0)
		room.world.positionPuck = NewVector(float64(GAME_WIDTH/2), float64(GAME_HEIGHT/2))

		if side == -1 {
			room.world.countA++
		} else if side == 1 {
			room.world.countB++
		}

		if room.world.countA >= MAX_GOALS || room.world.countB >= MAX_GOALS {
			pool.DeleteRoom(room.id, errors.New("game is finished"))
		}
	}
}

func (room *Room) broadcastWorldState() {
	puckPosition := NewPosition(int(math.Round(room.world.positionPuck.x)), int(math.Round(room.world.positionPuck.y)))

	room.playerA.GetWrite() <- NewWorldMessage(
		*room.world.positionA,
		*room.world.positionB,
		*puckPosition,
		room.world.countA,
		room.world.countB,
	)

	room.playerB.GetWrite() <- NewWorldMessage(
		*FlipPosition(room.world.positionB),
		*FlipPosition(room.world.positionA),
		*FlipPosition(puckPosition),
		room.world.countB,
		room.world.countA,
	)
}

func hitWall(magnitude *Vector, wallNormal *Vector) *Vector {
	mult := 2 * MultiplyVectors(magnitude, wallNormal)
	multNormal := MultiplyVectorNumber(wallNormal, mult)

	return SubstractVectors(magnitude, multNormal)
}

func detectPlayerHit(
	position *Position,
	prevPosition *Position,
	puck *Vector,
	puckPrev *Vector,
	puckMag *Vector,
	deltaT time.Duration,
) (bool, *Vector, *Vector) {
	lineStart := NewVector(float64(prevPosition.x), float64(prevPosition.y))
	lineEnd := NewVector(float64(position.x), float64(position.y))

	dX := lineEnd.x - lineStart.x
	dY := lineEnd.y - lineStart.y
	playerMagnitude := NewVector(dX/float64(deltaT.Milliseconds()), dY/float64(deltaT.Milliseconds()))

	if DistanceBetweenPoints(lineEnd, puck) < float64(2*ENTITY_RADIUS) {
		hitVector := NewVectorFromPoints(lineEnd, puck)
		magnitude := hitWall(puckMag, NormalizeVector(hitVector))
		resultMagnitude := SumVectors(magnitude, playerMagnitude)

		if (VectorLength(resultMagnitude)) < EPSILON {
			return false, nil, nil
		}

		positionAfterHit := PointOnLine(lineEnd, puck, float64(2*ENTITY_RADIUS)+COLLISION_DISTANCE)

		return true, positionAfterHit, resultMagnitude
	}

	collision, hitPoint := CheckSegmentSegmentIntercection(lineStart, lineEnd, puckPrev, puck)
	if collision {
		hitVector := NewVectorFromPoints(hitPoint, puck)
		magnitude := hitWall(puckMag, NormalizeVector(hitVector))
		resultMagnitude := SumVectors(magnitude, playerMagnitude)

		positionAfterHit := PointOnLine(hitPoint, puck, float64(2*ENTITY_RADIUS)+COLLISION_DISTANCE)

		return true, positionAfterHit, resultMagnitude
	}

	return false, nil, nil
}

func detectWallHit(walls [10][3]*Vector, puck *Vector, puckPrev *Vector) (bool, *Vector, *Vector) {
	for _, wall := range walls {
		collision, hitPoint := CheckSegmentSegmentIntercection(wall[0], wall[1], puck, puckPrev)
		if collision {
			hitAngle := AngleBetweenLines(puck, puckPrev, wall[0], wall[1])
			var distance float64
			angleDegree := RadToDegree(hitAngle)
			if math.Abs(angleDegree-90) < EPSILON {
				distance = float64(ENTITY_RADIUS)
			} else {
				distance = float64(ENTITY_RADIUS) / math.Sin(hitAngle)
			}

			hitPoint = PointOnLine(hitPoint, puckPrev, distance+COLLISION_DISTANCE)

			return true, hitPoint, wall[2]
		}

		collision, _ = CheckSegmentCircleIntercection(wall[0], wall[1], puck, float64(ENTITY_RADIUS))
		if collision {
			_, hitPoint := CheckLineLineIntercection(puck, puckPrev, wall[0], wall[1])
			hitAngle := AngleBetweenLines(puck, hitPoint, wall[0], wall[1])

			var distance float64
			angleDegree := RadToDegree(hitAngle)
			if math.Abs(angleDegree-90) < EPSILON {
				distance = float64(ENTITY_RADIUS)
			} else {
				distance = float64(ENTITY_RADIUS) / math.Sin(hitAngle)
			}

			hitPoint = PointOnLine(hitPoint, puck, distance+COLLISION_DISTANCE)

			return true, hitPoint, wall[2]
		}
	}

	return false, nil, nil
}

func detectGoal(puckPosition Vector) (bool, int) {
	if puckPosition.y-float64(ENTITY_RADIUS) < -float64(GATES_HEIGHT) {
		return true, -1
	}

	if puckPosition.y+float64(ENTITY_RADIUS) > float64(GAME_HEIGHT)+float64(GATES_HEIGHT) {
		return true, 1
	}

	return false, 0
}

func validateMagnitute(puckMag *Vector) *Vector {
	if VectorLength(puckMag) <= MAX_PUCK_MAGNITUDE {
		return puckMag
	}

	return ResizeVector(puckMag, MAX_PUCK_MAGNITUDE)
}
