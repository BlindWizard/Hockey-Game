package main

const EPSILON float64 = 0.0000001
const COLLISION_DISTANCE = 0.1
const DEFAULT_DRAG float64 = 0.01
const PLAYER_MESSAGE_THROTTLE = 20
const PHYSICS_CYCLE = 20
const NETWORK_CYCLE = 20

const GAME_WIDTH int = 800
const GAME_HEIGHT int = 1200
const GATES_WIDTH int = 300
const GATES_HEIGHT int = 80
const ENTITY_RADIUS int = 40

const MAX_GOALS uint = 10

type Position struct {
	x int
	y int
}

func NewPosition(x, y int) *Position {
	return &Position{x: x, y: y}
}

type World struct {
	positionA     *Position
	positionPrevA *Position
	positionB     *Position
	positionPrevB *Position
	positionPuck  *Vector
	magnitudePuck *Vector
	countA        uint
	countB        uint
	walls         [10][3]*Vector
}

func NewWorld() *World {
	return &World{
		positionA:     NewPosition(GAME_WIDTH/2, GAME_HEIGHT-2*ENTITY_RADIUS),
		positionPrevA: NewPosition(GAME_WIDTH/2, GAME_HEIGHT-2*ENTITY_RADIUS),
		positionB:     NewPosition(GAME_WIDTH/2, 2*ENTITY_RADIUS),
		positionPrevB: NewPosition(GAME_WIDTH/2, 2*ENTITY_RADIUS),
		positionPuck:  NewVector(float64(GAME_WIDTH/2), float64(GAME_HEIGHT/2)),
		magnitudePuck: NewVector(0, 0),
		countA:        0,
		countB:        0,
		walls:         buildWalls(),
	}
}

func buildWalls() [10][3]*Vector {
	return [10][3]*Vector{
		//top
		{NewVector(0, 0), NewVector(float64(GAME_WIDTH/2-GATES_WIDTH/2), 0), NewVector(0, 1)},
		{NewVector(float64(GAME_WIDTH/2+GATES_WIDTH/2), 0), NewVector(float64(GAME_WIDTH), 0), NewVector(0, 1)},
		//bot
		{NewVector(0, float64(GAME_HEIGHT)), NewVector(float64(GAME_WIDTH/2-GATES_WIDTH/2), float64(GAME_HEIGHT)), NewVector(0, -1)},
		{NewVector(float64(GAME_WIDTH/2+GATES_WIDTH/2), float64(GAME_HEIGHT)), NewVector(float64(GAME_WIDTH), float64(GAME_HEIGHT)), NewVector(0, -1)},
		//left
		{NewVector(0, 0), NewVector(0, float64(GAME_HEIGHT)), NewVector(1, 0)},
		//right
		{NewVector(float64(GAME_WIDTH), 0), NewVector(float64(GAME_WIDTH), float64(GAME_HEIGHT)), NewVector(-1, 0)},
		//top gate
		{NewVector(float64(GAME_WIDTH/2-GATES_WIDTH/2), 0), NewVector(float64(GAME_WIDTH/2-GATES_WIDTH/2), -float64(GATES_HEIGHT)), NewVector(1, 0)},
		{NewVector(float64(GAME_WIDTH/2+GATES_WIDTH/2), 0), NewVector(float64(GAME_WIDTH/2+GATES_WIDTH/2), -float64(GATES_HEIGHT)), NewVector(-1, 0)},
		//bottom gate
		{NewVector(float64(GAME_WIDTH/2-GATES_WIDTH/2), float64(GAME_HEIGHT)), NewVector(float64(GAME_WIDTH/2-GATES_WIDTH/2), float64(GAME_HEIGHT)+float64(GATES_HEIGHT)), NewVector(1, 0)},
		{NewVector(float64(GAME_WIDTH/2+GATES_WIDTH/2), float64(GAME_HEIGHT)), NewVector(float64(GAME_WIDTH/2+GATES_WIDTH/2), float64(GAME_HEIGHT)+float64(GATES_HEIGHT)), NewVector(-1, 0)},
	}
}

func FlipPosition(pos *Position) *Position {
	return &Position{
		x: GAME_WIDTH - pos.x,
		y: GAME_HEIGHT - pos.y,
	}
}

func validatePosition(pos *Position) *Position {
	pos.x = Clamp(pos.x, ENTITY_RADIUS, GAME_WIDTH-ENTITY_RADIUS)
	pos.y = Clamp(pos.y, GAME_HEIGHT/2+ENTITY_RADIUS, GAME_HEIGHT-ENTITY_RADIUS)

	return pos
}
