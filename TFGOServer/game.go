package TFGOServer

import (
	"time"
)

type Location struct {
	X float64
	Y float64
}

type Game struct {
	ID string
	Name string
	Password string

	Status int
	Mode int
	Timer time.Timer

	RedTeam Team
	RedBase Location
	BlueTeam Team
	BlueBase Location

	Boundaries []Location // ul, ur, lr, ll
	ControlPoints map[string]ControlPoint
}
