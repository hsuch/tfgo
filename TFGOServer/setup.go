package main

// setup.go: functions for game setup

import (
	"time"
	"math"
	"math/rand"
	"net"
	"encoding/json"
	"strconv"
	"fmt"
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

	fmt.Printf("Player %s created.\n", p.Name)

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
		p := Location{X: degreeToMeter(vertex["X"].(float64)), Y: degreeToMeter(vertex["Y"].(float64))}
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

	fmt.Printf("Game %s with ID %s created.\n", g.Name, g.ID)
	fmt.Printf("Player %s added to game %s.\n", p.Name, g.ID)

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

		fmt.Printf("Player %s added to game %s.\n", p.Name, gameID)

		return target
	}
}

// determine locations and radii of bases and control points
func (g *Game) generateObjectives(numCP int) {
	minX := math.MaxFloat64
	maxX := -math.MaxFloat64
	minY := math.MaxFloat64
	maxY := -math.MaxFloat64
	for _, val := range g.Boundaries {
		if val.P.X < minX {
			minX = val.P.X
		}
		if val.P.X > maxX {
			maxX = val.P.X
		}
		if val.P.Y < minY {
			minY = val.P.Y
		}
		if val.P.Y > maxY {
			maxY = val.P.Y
		}
	}
	xrange := maxX - minX
	yrange := maxY - minY

	// set up base locations for the two teams
	baseRadius := BASERADIUS()
	g.RedTeam.BaseRadius = baseRadius
	g.BlueTeam.BaseRadius = baseRadius
	offset := baseRadius + 2
	if xrange > yrange {
		mid := yrange / 2
		g.RedTeam.Base = Location{maxX - offset, mid}
		g.BlueTeam.Base = Location{minX + offset, mid}
	} else {
		mid := xrange / 2
		g.RedTeam.Base = Location{mid, maxY - offset}
		g.BlueTeam.Base = Location{mid, minY + offset}
	}

	// set up control points
	cpRadius := CPRADIUS()
	// make sure that control points don't intersect bases
	minX = minX + 2 * offset + cpRadius
	maxX = maxX - 2 * offset - cpRadius
	minY = minY + 2 * offset + cpRadius
	maxY = maxY - 2 * offset - cpRadius
	xrange = maxX - minX
	yrange = maxY - minY
	for i := 0; i < numCP; i++ {
		cpLoc := Location{minX + r.Float64() * xrange, minY + r.Float64() * yrange}
		if inBounds(g, cpLoc) {
			id := "CP" + strconv.Itoa(i)
			cp := &ControlPoint{ID: id, Location: cpLoc, Radius: cpRadius}
			g.ControlPoints[id] = cp
		} else {
			i-- // if this location is invalid, decrement i so that it doesn't count towards numCP
		}
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

// begin a game, determining objective and team information and
// starting goroutines that will run for the duration of the game
func (g *Game) start() {
	// for now we're doing just one ControlPoint, that may change later
	g.generateObjectives(1)
	g.randomizeTeams()

	startTime := time.Now().Add(time.Minute)
	sendGameStartInfo(g, startTime)
	go sendGameUpdates(g)

	time.Sleep(time.Until(startTime))
	g.Status = PLAYING
	g.Timer = time.AfterFunc(g.TimeLimit, func() {
		g.stop()
	})
	for _, cp := range g.ControlPoints {
		go cp.updateStatus(g)
	}
}

// end a game, signalling and performing resource cleanup
func (g *Game) stop() {
	sendGameOver(g)
	g.Status = GAMEOVER
	delete(games, g.ID)
	for _, player := range g.Players {
		close(player.Chan)
	}
}
