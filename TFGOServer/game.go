package main

import (
	"time"
	"math"
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

func (l1 Location) getDistance(l2 Location) float64 {
	first := math.Pow(float64(l2.X-l1.X), 2)
	second := math.Pow(float64(l2.Y-l1.Y), 2)
	return math.Sqrt(first + second)
}