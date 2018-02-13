package main

import (
	"time"
	"math"
	"math/rand"
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

// the following random string generation code is heavily inspired by the
// example code at https://siongui.github.io/2015/04/13/go-generate-random-string/
var r = rand.New(rand.NewSource(time.Now().UnixNano()))
const idChars = "abcdefghijklmnopqrstuvwxyz1234567890"

func createGameID() string {
	result := make([]byte, 16)
	for i := range result {
		result[i] = idChars[r.Intn(len(idChars))]
	}
	return string(result)
}

func createNewGame(data map[string]interface{}) (*Game, *Player) {
	var g Game
	g.ID = createGameID()
	g.Name = data["Name"].(string)
	g.Password = data["Password"].(string)
	g.Description = data["Description"].(string)
	g.PlayerLimit = int(data["PlayerLimit"].(float64))
	g.PointLimit = int(data["PointLimit"].(float64))
	g.TimeLimit, _ = time.ParseDuration(data["TimeLimit"].(string))
	g.Status = CREATING
	g.Mode = stringToMode[data["Mode"].(string)]

	boundaries := data["Boundaries"].([]interface{})
	for _, val := range boundaries {
		vertex := val.(map[string]interface{})
		p := Location{X: vertex["X"].(float64), Y: vertex["Y"].(float64)}
		// calculate d, t here
		border := Border{P: p}
		g.Boundaries = append(g.Boundaries, border)
	}


	return &g
}