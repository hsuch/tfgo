package main

// setup.go: functions for game setup

import (
	"time"
	"math/rand"
	"net"
	"encoding/json"
)

func (g *Game) findCenter() Location {
	X := 0.0
	Y := 0.0
	c := float64(len(g.Boundaries))
	for _, val := range g.Boundaries {
		X += val.P.X
		Y += val.P.Y
	}
	return Location{X / c, Y / c}
}

func createPlayer(conn net.Conn, name, icon string) *Player {
	var p Player
	p.Name = name
	p.Icon = icon
	p.Team = NEUTRAL

	p.Conn = conn
	p.Chan = make(chan map[string]interface{})
	p.Encoder = json.NewEncoder(conn)

	p.Status = NORMAL
	p.Health = 100
	p.Armor = 0

	go p.sender()

	return &p
}

// the following random string generation code is heavily inspired by the
// example code at https://siongui.github.io/2015/04/13/go-generate-random-string/
var r = rand.New(rand.NewSource(time.Now().UnixNano()))
const idChars = "abcdefghijklmnopqrstuvwxyz1234567890"

func createGameID() string {
	candidate := make([]byte, 16)
	for i := range candidate {
		candidate[i] = idChars[r.Intn(len(idChars))]
	}

	if _, exists := games[string(candidate)]; exists {
		return createGameID()
	} else {
		return string(candidate)
	}
}

func createGame(conn net.Conn, data map[string]interface{}) (*Game, *Player) {
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
		g.Boundaries = append(g.Boundaries, Border{P: p})
	}

	for i, boundary := range g.Boundaries {
		var index int
		if i == 0 {
			index = len(g.Boundaries) - 1
		} else {
			index = i - 1
		}

		prev := g.Boundaries[index].P
		g.Boundaries[index].D = Direction{boundary.P.X - prev.X, boundary.P.Y - prev.Y}
	}

	host := data["Host"].(map[string]interface{})
	p := createPlayer(conn, host["Name"].(string), host["Icon"].(string))

	// during game creation, everyone is placed on the red team as NEUTRAL.
	// starting the game will randomly assign teams
	g.RedTeam = &Team{Name: "Red Team", Players: map[string]*Player{p.Name : p}}
	g.BlueTeam = &Team{Name: "Blue Team", Players: map[string]*Player{}}

	games[g.ID] = &g
	return &g, p
}

func (p *Player) joinGame(gameID string) *Game {
	target := games[gameID]
	if len(target.RedTeam.Players) == target.PlayerLimit {
		sendJoinGameError(p, "GameFull")
		return nil
	} else if target.Status != CREATING {
		sendJoinGameError(p, "GameStarted")
		return nil
	} else {
		target.RedTeam.Players[p.Name] = p
		sendPlayerListUpdate(target)
		return target
	}
}

func (g *Game) randomizeTeams() {
	teamSize := len(g.RedTeam.Players) / 2
	count := 0

	// iteration order through maps is random
	for _, player := range g.RedTeam.Players {
		if count < teamSize {
			delete(g.RedTeam.Players, player.Name)
			g.BlueTeam.Players[player.Name] = player
			count++
		} else {
			break
		}
	}

	for _, player := range g.RedTeam.Players {
		player.Team = RED
	}
	for _, player := range g.BlueTeam.Players {
		player.Team = BLUE
	}
}

func (g *Game) generateObjectives() {
	// jenny halp
	// determine location and radius of bases and control points
	// teams have already been defined, so simply set the Base, BaseRadius fields
	// control points have not been defined, so the control point must be created,
	// all relevant fields filled out, and its pointer set in the map of the game's
	// ControlPoints
}

func (g *Game) start() {
	g.randomizeTeams()
	g.generateObjectives()

	startTime := time.Now().Add(time.Minute)
	sendGameStartInfo(g, startTime)

	time.Sleep(time.Until(startTime))
	g.Status = PLAYING
	g.Timer = time.NewTimer(g.TimeLimit)
	go g.awaitGameEnd()
	go sendGameUpdates(g)
	for _, cp := range g.ControlPoints {
		go cp.updateStatus(g)
	}
}