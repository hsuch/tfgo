package main

import (
	"time"
	"net"
	"encoding/json"
	"math"
)

// definitions and constants for game structures

type Allegiance int
const (
	RED Allegiance = iota
	NEUTRAL
	BLUE
)

type PlayerStatus int
const (
	NORMAL PlayerStatus = iota
	OUTOFBOUNDS
	RESPAWNING
)

var playerStatusToString = map[PlayerStatus]string {
	NORMAL : "NORMAL",
	OUTOFBOUNDS : "OUTOFBOUNDS",
	RESPAWNING : "RESPAWNING",
}

type Player struct {
	Name string
	Icon string
	Conn net.Conn
	Chan chan map[string]interface{}
	Encoder *json.Encoder
	Team Allegiance
	OccupyingPoint *ControlPoint

	Status PlayerStatus
	StatusTimer *time.Timer

	Health int
	Armor int

	Inventory map[string]Pickup

	Location Location
}

type Team struct {
	Name string
	Players map[string]*Player
	Base Location
	Points int
}

type ControlPoint struct {
	ID string

	Location Location
	Radius float64

	PayloadPath [2]Location // start, end
	PayloadLoc Location

	// number of currently occupying players from each team
	RedCount int
	BlueCount int

	// number in [-100, 100] indicating capture/decapture progress.
	// hitting -100 or 100 from neutral ownership yields control
	// control point to red or blue, respectively. hitting 0
	// from team ownership neutralizes control point.
	CaptureProgress int
	ControllingTeam Allegiance
}

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

var modeToString = map[Mode]string {
	SINGLECAP : "SingleCapture",
	MULTICAP : "MultiCapture",
	PAYLOAD : "Payload",
}

var stringToMode = map[string]Mode {
	"SingleCapture" : SINGLECAP,
	"MultiCapture" : MULTICAP,
	"Payload" : PAYLOAD,
}

type Location struct {
	X float64
	Y float64
}

type Direction struct {
	X float64
	Y float64
}

type Border struct {
	P Location	// one of the two vertices which define this border
	D Direction	// the direction vector of the line
	T float64	// the max t-value for this line segment
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

var games = make(map[string]*Game)

type Pickup interface {
	use(game *Game, player *Player)
}

// each of the available weapons is defined as a globally
// accessible variable
var SWORD = Weapon {
	Name: "Sword",
	Damage: 25,
	Spread: 2*math.Pi,
	Range: 1,
	ClipSize: 500,
	ShotReload: time.Second * 0,
	ClipReload: time.Second * 0,
}

var SHOTGUN = Weapon {
	Name: "Shotgun",
	Damage: 25,
	Spread: math.Pi/2,
	Range: 3,
	ClipSize: 2,
	ShotReload: time.Millisecond * 500,
	ClipReload: time.Second * 3,
}

var weapons = map[string]Weapon {
	"Sword" : SWORD,
	"Shotgun" : SHOTGUN,
}

type Weapon struct {
	Name string

	Damage int
	Spread float64
	Range float64

	ClipSize int
	ShotReload time.Duration
	ClipReload time.Duration
}