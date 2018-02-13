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

type Game struct {
	ID string
	Name string
	Password string

	Status GameStatus
	Mode Mode
	Timer *time.Timer

	RedTeam *Team
	BlueTeam *Team

	Boundaries [4]Location // ul, ur, lr, ll
	ControlPoints map[string]*ControlPoint
}

func (l1 Location) getDistance(l2 Location) float64 {
	first := math.Pow(float64(l2.X-l1.X), 2)
	second := math.Pow(float64(l2.Y-l1.Y), 2)
	return math.Sqrt(first + second)
}