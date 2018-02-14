package main

import (
	"time"
	"math/rand"
	"net"
	"encoding/json"
)

func (game *Game) findCenter() Location {
	X := 0.0
	Y := 0.0
	c := len(game.Boundaries)
	_, val := range game.Boundaries {
		X += val.P.X
		Y += val.P.Y
	}
	return Location{X / c, Y / c}
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
	for key, val := range boundaries {
		vertex := val.(map[string]interface{})
		p := Location{X: vertex["X"].(float64), Y: vertex["Y"].(float64)}
		border := Border{P: p}
		prev := g.Boundaries[key - 1].P
		g.Boundaries[key - 1].D = Direction{p.X - prev.X, p.Y - prev.Y}
		g.Boundaries = append(g.Boundaries, border)
	}
	first := g.Boundaries[0].P
	last := g.Boundaries[len(g.Boundaries) - 1].P
	g.Boundaries[len(g.Boundaries) - 1].D = Direction{first.X - last.X, first.Y - last.Y}

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
