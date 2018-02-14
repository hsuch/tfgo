package main

// setup.go: functions for game setup

import (
	"time"
	"math/rand"
	"net"
	"encoding/json"
)

// return the central position of a game
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

// register a new player, but do not add them to any team.
// performed whenever someone chooses Create Game or Join
// Game from the primary app menu
func createPlayer(conn net.Conn, name, icon string) *Player {
	var p Player
	p.Name = name
	p.Icon = icon

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

// generate a random unique game ID
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

// determine game boundaries based on vertex information
func (g *Game) setBoundaries(boundaries []interface{}) {
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
}

// register a new game instance, with the host as its first player
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
	g.setBoundaries(data["Boundaries"].([]interface{}))

	host := data["Host"].(map[string]interface{})
	p := createPlayer(conn, host["Name"].(string), host["Icon"].(string))

	g.RedTeam = &Team{Name: "Red"}
	g.BlueTeam = &Team{Name: "Blue"}
	g.Players = map[string]*Player{p.Name : p}

	games[g.ID] = &g
	return &g, p
}

// add a player to a game if possible
func (p *Player) joinGame(gameID string) *Game {
	target := games[gameID]
	if len(target.Players) == target.PlayerLimit {
		sendJoinGameError(p, "GameFull")
		return nil
	} else if target.Status != CREATING {
		sendJoinGameError(p, "GameStarted")
		return nil
	} else {
		target.Players[p.Name] = p
		sendPlayerListUpdate(target)
		return target
	}
}

// assign players to teams at the start of a game
func (g *Game) randomizeTeams() {
	teamSize := len(g.Players) / 2
	count := 0

	// iteration order through maps is random
	for _, player := range g.Players {
		if count < teamSize {
			player.Team = g.RedTeam
		} else {
			player.Team = g.BlueTeam
		}
		count++
	}
}

// determine locations and radii of bases and control points
func (g *Game) generateObjectives() {
	// jenny halp
	// determine location and radius of bases and control points
	// teams have already been defined, so simply set the Base, BaseRadius fields
	// control points have not been defined, so the control point must be created,
	// all relevant fields filled out, and its pointer set in the map of the game's
	// ControlPoints
}

// begin a game, determining objective and team information and
// starting goroutines that will run for the duration of the game
func (g *Game) start() {
	g.generateObjectives()
	g.randomizeTeams()

	startTime := time.Now().Add(time.Minute)
	sendGameStartInfo(g, startTime)

	time.Sleep(time.Until(startTime))
	g.Status = PLAYING
	g.Timer = time.AfterFunc(g.TimeLimit, func() {
		g.stop()
	})
	go sendGameUpdates(g)
	for _, cp := range g.ControlPoints {
		go cp.updateStatus(g)
	}
}

// end a game, signalling and performing resource cleanup
func (g *Game) stop() {
	g.Status = GAMEOVER
	delete(games, g.ID)
	for _, player := range g.Players {
		close(player.Chan)
	}
}
