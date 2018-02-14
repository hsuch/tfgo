package main

// setup.go: functions for game setup

import (
	"time"
	"math/rand"
	"net"
	"encoding/json"
)

func (game *Game) findCenter() Location {
	// jenny halp
	return Location{}
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

func createNewGame(conn net.Conn, data map[string]interface{}) (*Game, *Player) {
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

	var p Player
	host := data["Host"].(map[string]interface{})
	p.Name = host["Name"].(string)
	p.Icon = host["Icon"].(string)
	p.Conn = conn
	p.Chan = make(chan map[string]interface{})
	p.Encoder = json.NewEncoder(conn)
	p.Team = NEUTRAL
	p.respawn(&g)
	go p.sender()

	// during game creation, everyone is placed on the red team as NEUTRAL.
	// starting the game will randomly assign teams
	g.RedTeam = &Team{Name: "Red Team", Players: map[string]*Player{p.Name : &p}}
	g.BlueTeam = &Team{Name: "Blue Team", Players: map[string]*Player{}}

	games[g.ID] = &g

	sendPlayerListUpdate(&g)

	return &g, &p
}