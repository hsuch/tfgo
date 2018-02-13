package main

import (
	"time"
)

type GameStatus int
const (
	CREATING GameStatus = iota
	PLAYING
	GAMEOVER
)

type Mode int
const (
	SINGLECAP Mode = iota
	MULTICAP
	PAYLOAD
)

type Location struct {
	X float64
	Y float64
}

type Border struct {
	P Location	//one of the two vertices which define this border
	D Direction	//the direction vector of the line
	T float64	//the max t-value for this line segment
}

type Game struct {
	ID string
	Name string
	Password string
	Description string

	PlayerLimit int
	PointLimit int
	TimeLimit time.Duration

	Status GameStatus
	Mode Mode
	Timer *time.Timer

	RedTeam *Team
	BlueTeam *Team

	Boundaries []Border
	ControlPoints map[string]*ControlPoint
}
