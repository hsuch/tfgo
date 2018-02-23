package main

// struct.go: definitions and constants for game structures

import (
	"time"
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
	Players  map[string]*Player

	Boundaries    []Border
	ControlPoints map[string]*ControlPoint
	Pickups		  []PickupSpot
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

	Chan chan map[string]interface{} // used to synchronize sends
	Encoder *json.Encoder

	Status PlayerStatus
	StatusTimer *time.Timer // tracks duration of abnormal statuses

	Health int
	Armor int
	SelectedWeapon Weapon
	Weapons map[string]Weapon
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

	// number in [-7, 7] indicating capture/decapture progress.
	// hitting -7 or 7 from neutral ownership yields control
	// control point to red or blue, respectively. hitting 0
	// from team ownership neutralizes control point.
	CaptureProgress int
	ControllingTeam *Team
}

type PickupSpot struct {
	Location Location
	Pickup Pickup
	Available bool
	SpawnTimer *time.Timer // duration until respawn
}

// since pickups can vary wildly, we use an interface rather than
// a type and only require them to implement a use() method
type Pickup interface {
	use(player *Player)
}

type ArmorPickup struct {
	AP int
}

type HealthPickup struct {
	HP int
}

type WeaponPickup struct {
	WP Weapon
}

type Weapon struct {
	Name string

	Damage int
	Spread float64
	Range  float64

	ClipSize   int
	ShotReload time.Duration
	ClipReload time.Duration
}

var weapons = map[string]Weapon {
	"Sword" : SWORD,
	"Shotgun" : SHOTGUN,
}

var weaponsSlice = []Weapon {
	SWORD,
	SHOTGUN,
}

// each of the available weapons is defined as a globally
// accessible variable
var SWORD = Weapon {
	Name: "Sword",
	Damage: 25,
	Spread: math.Pi,
	Range: 50,
	ClipSize: 1337,
	ShotReload: time.Second * 0,
	ClipReload: time.Second * 0,
}

var SHOTGUN = Weapon {
	Name: "Shotgun",
	Damage: 25,
	Spread: math.Pi/2,
	Range: 5,
	ClipSize: 2,
	ShotReload: time.Millisecond * 500,
	ClipReload: time.Second * 3,
}

// Helper functions, mostly for conversions
func meterToDegree(m float64) float64 {
	return m * 9 / 1000000
}

func degreeToMeter(d float64) float64 {
	return d * 1000000 / 9
}

func (l Location) locationToDegrees() Location {
	return Location{X: meterToDegree(l.X), Y: meterToDegree(l.Y)}
}

func intMin(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func intMax(a, b int) int {
	if a <= b {
		return b
	}
	return a
}

// constants defined via functions, as Go does not allow
// for non-primitive constants
func TICK() time.Duration {
	return 200 * time.Millisecond
}

func OUTOFBOUNDSTIME() time.Duration {
	return 10 * time.Millisecond
}

func RESPAWNTIME() time.Duration {
	return 15 * time.Second
}

func PICKUPRESPAWNTIME() time.Duration {
	return 15 * time.Second
}


// returns the baseRadius given the games x and y dimensions
// default is 3m, but size is adjusted down if dimensions are too small
func BASERADIUS(x, y float64) float64 {
	if x < 14 || y < 14 {
		return math.Min(x, y) * 5 / 14
	} else if x < 28 && y < 28 {
		return math.Max(x, y) * 5 / 28
	} else {
		return 5.0
	}
}

func CPRADIUS() float64 {
	return 3.0
}

func MAXHEALTH() int {
	return 100
}

func MAXARMOR() int {
	return 100
}

func PICKUPRADIUS() float64 {
	return 1.0
}

func PICKUPDISTRIBUTION() float64 {
	return 10.0
}
