package main

import (
	"time"
	"reflect"
)

// clientmessage.go: functions for building and sending messages to clients

// goroutine responsible for delivering messages to players
func (p *Player) sender() {
	for {
		msg, closed := <-p.Chan
		if closed {
			return
		} else {
			p.Encoder.Encode(msg)
		}
	}
}

// deliver a message to all players
func (g *Game) broadcast(msg map[string]interface{}) {
	for _, player := range g.Players {
		player.Chan <- msg
	}
}


// utility functions that help construct various message components

// this method of retrieving field values by name inspired by
// https://stackoverflow.com/questions/18930910/golang-access-struct-property-by-name
func (g *Game) getPlayerInfo(fields []string) []map[string]interface{} {
	var playerList []map[string]interface{}
	for _, player := range g.Players {
		playerInfo := make(map[string]interface{})
		r := reflect.ValueOf(player)
		for _, field := range fields {
			if field == "Team" {
				playerInfo[field] = player.Team.Name
			} else {
				f := reflect.Indirect(r).FieldByName(field)
				playerInfo[field] = f.Interface()
			}
		}
		playerList = append(playerList, playerInfo)
	}

	return playerList
}

func (g *Game) getBoundaryVertices() []map[string]float64 {
	var vertices []map[string]float64
	for _, boundary := range g.Boundaries {
		vertex := make(map[string]float64)
		vertex["X"] = boundary.P.X
		vertex["Y"] = boundary.P.Y
		vertices = append(vertices, vertex)
	}

	return vertices
}

func (t *Team) getLocInfo() map[string]interface{} {
	return map[string]interface{} {
		"Location" : t.Base,
		"Radius" : t.BaseRadius,
	}
}

func (cp *ControlPoint) getLocInfo() map[string]interface{} {
	return map[string]interface{} {
		"Location" : cp.Location,
		"Radius" : cp.Radius,
	}
}

func (g *Game) getObjectiveUpdate() []map[string]interface{} {
	occupants := make(map[*ControlPoint][]string)
	for _, player := range g.Players {
		cp := player.OccupyingPoint
		if cp != nil {
			occupants[cp] = append(occupants[cp], player.Name)
		}
	}

	var cpList []map[string]interface{}
	for _, cp := range g.ControlPoints {
		cpInfo := make(map[string]interface{})
		cpInfo["Location"] = cp.Location
		cpInfo["Occupying"] = occupants[cp]
		if cp.ControllingTeam == nil {
			cpInfo["BelongsTo"] = "Neutral"
		} else {
			cpInfo["BelongsTo"] = cp.ControllingTeam.Name
		}
		cpInfo["Progress"] = cp.CaptureProgress
		cpList = append(cpList, cpInfo)
	}

	return cpList
}

// each sendX function corresponds to the server to client message with
// "Action": X in https://github.com/hsuch/tfgo/wiki/Network-Messages

func sendPlayerListUpdate(game *Game) {
	playerList := game.getPlayerInfo([]string{"Name", "Icon"})
	msg := map[string]interface{} {
		"Type" : "PlayerListUpdate",
		"Data" : playerList,
	}
	game.broadcast(msg)
}

func sendAvailableGames(player *Player) {
	var gameList []map[string]interface{}
	for _, game := range games {
		if game.Status == CREATING {
			gameInfo := make(map[string]interface{})
			gameInfo["ID"] = game.ID
			gameInfo["Name"] = game.Name
			gameInfo["Mode"] = modeToString[game.Mode]
			gameInfo["Location"] = game.findCenter()
			gameInfo["PlayerList"] = game.getPlayerInfo([]string{"Name", "Icon"})
			gameList = append(gameList, gameInfo)
		}
	}

	msg := map[string]interface{} {
		"Type" : "AvailableGames",
		"Data" : gameList,
	}
	player.Chan <- msg
}

func sendGameInfo(player *Player, gameID string) {
	target := games[gameID]
	gameInfo := make(map[string]interface{})
	gameInfo["Description"] = target.Description
	gameInfo["PlayerLimit"] = target.PlayerLimit
	gameInfo["PointLimit"] = target.PointLimit
	gameInfo["TimeLimit"] = target.TimeLimit.String()
	gameInfo["Boundaries"] = target.getBoundaryVertices()
	gameInfo["PlayerList"] = target.getPlayerInfo([]string{"Name", "Icon"})

	msg := map[string]interface{} {
		"Type" : "GameInfo",
		"Data" : gameInfo,
	}
	player.Chan <- msg
}

func sendJoinGameError(player *Player, error string) {
	msg := map[string]interface{} {
		"Type" : "JoinGameError",
		"Data" : error,
	}
	player.Chan <- msg
}

func sendGameStartInfo(game *Game, startTime time.Time) {
	gameInfo := make(map[string]interface{})
	gameInfo["PlayerList"] = game.getPlayerInfo([]string{"Name", "Team"})
	gameInfo["RedBase"] = game.RedTeam.getLocInfo()
	gameInfo["BlueBase"] = game.BlueTeam.getLocInfo()
	var cpInfo []map[string]interface{}
	for _, cp := range game.ControlPoints {
		cpInfo = append(cpInfo, cp.getLocInfo())
	}
	gameInfo["Objectives"] = cpInfo
	gameInfo["StartTime"] = startTime.Format("2006-01-02 15:04:05")

	msg := map[string]interface{} {
		"Type" : "GameStartInfo",
		"Data" : gameInfo,
	}
	game.broadcast(msg)
}

func sendGameUpdates(game *Game) {
	for game.Status == PLAYING {
		gameInfo := make(map[string]interface{})
		gameInfo["PlayerList"] = game.getPlayerInfo([]string{"Name", "Orientation", "Location"})
		gameInfo["Points"] = map[string]int {
			"Red" : game.RedTeam.Points,
			"Blue" : game.BlueTeam.Points,
		}
		gameInfo["Objectives"] = game.getObjectiveUpdate()

		msg := map[string]interface{} {
			"Type" : "GameUpdate",
			"Data" : gameInfo,
		}
		game.broadcast(msg)

		time.Sleep(TICK())
	}
}
