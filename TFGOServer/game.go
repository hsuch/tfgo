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
