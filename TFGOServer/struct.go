package main

// struct.go: definitions and constants for game structures

import (
	"time"
	"net"
	"encoding/json"
	"math"
)

// global collection of all active games
var games = make(map[string]*Game)

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

// mode/string maps to facilitate JSON communication
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

// represents a point in space
type Location struct {
	X float64
	Y float64
}

// represents a direction vector
type Direction struct {
	X float64
	Y float64
}

// represents a bounding edge for the game arena
type Border struct {
	P Location	// one of the two vertices which define this border
	D Direction	// the direction vector of the line
}

type Game struct {
	ID          string
	Name        string
	Password    string
	Description string

	PlayerLimit int
	PointLimit  int
	TimeLimit   time.Duration

	Status GameStatus
	Mode   Mode
	Timer  *time.Timer

	RedTeam  *Team
	BlueTeam *Team
	Players map[string]*Player

	Boundaries    []Border
	ControlPoints map[string]*ControlPoint
}

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
	Team *Team

	Conn net.Conn
	Chan chan map[string]interface{} // used to synchronize sends
	Encoder *json.Encoder

	Status PlayerStatus
	StatusTimer *time.Timer // tracks duration of abnormal statuses

	Health int
	Armor int
	Inventory map[string]Pickup
	Location Location
	Orientation float64
	OccupyingPoint *ControlPoint // control point player is currently in
}

type Team struct {
	Name string
	Points int

	Base Location
	BaseRadius float64
}

type ControlPoint struct {
	ID string

	Location Location
	Radius float64

	// only used for payload games
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
	ControllingTeam *Team
}

// since pickups can vary wildly, we use an interface rather than
// a type and only require them to implement a use() method
type Pickup interface {
	use(game *Game, player *Player)
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

var weapons = map[string]Weapon {
	"Sword" : SWORD,
	"Shotgun" : SHOTGUN,
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

// constants defined via functions, as Go does not allow
// for non-primitive constants
func TICK() time.Duration {
	return 200 * time.Millisecond
}

func meterToDegree(m float64) float64 {
	return m * 9 / 1000000
}

func degreeToMeter(d float64) float64 {
	return d * 1000000 / 9
}

func BASERADIUS() float64 {
	return meterToDegree(3.0)
}

func CPRADIUS() float64 {
	return meterToDegree(1.0)
}
