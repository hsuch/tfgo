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
	HostID      string
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
	Pickups		  []*PickupSpot

	PayloadPath		Direction	//direction of Payload motion from Red base to Blue base
	PayloadSpeed	float64		//in meter/second
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
	ID string
	Name string
	Icon string
	Team *Team

	Chan chan map[string]interface{} // used to synchronize sends
	Encoder *json.Encoder

	Status PlayerStatus
	StatusTimer *time.Timer // tracks duration of abnormal statuses

	Health int
	Armor int
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
	"Pistol" : PISTOL,
	"Blaster" : BLASTER,
	"Crossbow" : CROSSBOW,
	"SniperRifle" : RIFLE,
	"Boomerang" : BOOMERANG,
	"Lightsaber" : LIGHTSABER,
	"Spear" : SPEAR,
	"BanHammer" : BANHAMMER,
	"BeeSwarm" : BEESWARM,
}

var weaponToString = map[Weapon]string {
	SWORD : "Sword",
	SHOTGUN : "Shotgun",
	PISTOL: "Pistol",
	BLASTER: "Blaster",
	CROSSBOW: "Crossbow",
	RIFLE: "SniperRifle",
	BOOMERANG: "Boomerang",
	LIGHTSABER: "Lightsaber",
	SPEAR: "Spear",
	BANHAMMER: "BanHammer",
	BEESWARM: "BeeSwarm",
}

// each of the available weapons is defined as a globally
// accessible variable
var SWORD = Weapon {
	Name: "Sword",
	Damage: 25,
	Spread: math.Pi,
	Range: 3,
	ClipSize: 1337,
	ShotReload: time.Second * 0,
	ClipReload: time.Second * 0,
}

var SHOTGUN = Weapon {
	Name: "Shotgun",
	Damage: 25,
	Spread: math.Pi/2,
	Range: 12,
	ClipSize: 2,
	ShotReload: time.Millisecond * 500,
	ClipReload: time.Second * 3,
}

var PISTOL = Weapon {
	Name: "Pistol",
	Damage: 25,
	Spread: math.Pi/2,
	Range: 10,
	ClipSize: 6,
	ShotReload: time.Millisecond * 500,
	ClipReload: time.Second * 3,
}

var BLASTER = Weapon {
	Name: "Blaster",
	Damage: 35,
	Spread: math.Pi/3,
	Range: 16,
	ClipSize: 10,
	ShotReload: time.Millisecond * 300,
	ClipReload: time.Second * 5,
}

var CROSSBOW = Weapon {
	Name: "Crossbow",
	Damage: 20,
	Spread: math.Pi/2,
	Range: 10,
	ClipSize: 20,
	ShotReload: time.Second * 3,
	ClipReload: time.Second * 20,
}

var RIFLE = Weapon {
	Name: "SniperRifle",
	Damage: 30,
	Spread: math.Pi/8,
	Range: 40,
	ClipSize: 10,
	ShotReload: time.Millisecond * 300,
	ClipReload: time.Second * 3,
}

var BOOMERANG = Weapon {
	Name: "Boomerang",
	Damage: 20,
	Spread: math.Pi,
	Range: 13,
	ClipSize: 1,
	ShotReload: time.Second * 0,
	ClipReload: time.Second * 20,
}

var LIGHTSABER = Weapon {
	Name: "Lightsaber",
	Damage: 50,
	Spread: math.Pi,
	Range: 3,
	ClipSize: 1337,
	ShotReload: time.Second * 0,
	ClipReload: time.Second * 0,
}

var SPEAR = Weapon {
	Name: "Spear",
	Damage: 20,
	Spread: math.Pi/2,
	Range: 8,
	ClipSize: 3,
	ShotReload: time.Second * 1,
	ClipReload: time.Second * 20,
}

var BANHAMMER = Weapon {
	Name: "BanHammer",
	Damage: 30,
	Spread: math.Pi/2,
	Range: 15,
	ClipSize: 1,
	ShotReload: time.Second * 0,
	ClipReload: time.Second * 10,
}

var BEESWARM = Weapon {
	Name: "BeeSwarm",
	Damage: 10,
	Spread: 2 * math.Pi,
	Range: 50,
	ClipSize: 1337,
	ShotReload: time.Second * 0,
	ClipReload: time.Second * 0,
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
	return 10 * time.Second
}

func RESPAWNTIME() time.Duration {
	return 15 * time.Second
}

func PICKUPRESPAWNTIME() time.Duration {
	return 15 * time.Second
}

// returns the baseRadius given the game's x and y dimensions
// default is 7m, but size is adjusted down if dimensions are too small
func BASERADIUS(x, y float64) float64 {
	minLen := 4.0 + 7.0 * 2.0
	if x < minLen || y < minLen {
		return math.Min(x, y) * 7 / minLen
	} else if x < (minLen * 2) && y < (minLen * 2) {
		return math.Max(x, y) * 7 / (minLen * 2)
	} else {
		return 7.0
	}
}

func CPRADIUS() float64 {
	return 5.0
}

func PICKUPRADIUS() float64 {
	return 3.0
}

// calculates the pickup distribution, which is either
// one every 10 m^2 or whatever distribution will result in a
// maximum of 25 pickups per game
func PICKUPDISTRIBUTION(xRange, yRange float64) float64 {
	dist := math.Sqrt((xRange * yRange) / 25)
	return math.Max(dist, 10.0)
}

func MAXHEALTH() int {
	return 100
}

func MAXARMOR() int {
	return 100
}

func MAXSPEED() float64 {
	return 0.5
}
