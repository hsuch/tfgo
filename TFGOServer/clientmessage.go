package main

import "time"

// clientmessage.go: functions for building and sending messages to clients

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

func (g *Game) getLobbyList() []map[string]string {
	var playerList []map[string]string
	for _, player := range g.RedTeam.Players {
		playerInfo := make(map[string]string)
		playerInfo["Name"] = player.Name
		playerInfo["Icon"] = player.Icon
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

func (g *Game) getTeamList() []map[string]string {
	var teamList []map[string]string
	for _, player := range g.RedTeam.Players {
		playerInfo := make(map[string]string)
		playerInfo["Name"] = player.Name
		playerInfo["Team"] = allegianceToString[player.Team]
		teamList = append(teamList, playerInfo)
	}

	return teamList
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

func (g *Game) broadcast(msg map[string]interface{}) {
	for _, player := range g.RedTeam.Players {
		player.Chan <- msg
	}
	for _, player := range g.BlueTeam.Players {
		player.Chan <- msg
	}
}

func sendPlayerListUpdate(game *Game) {
	playerList := game.getLobbyList()
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
			gameInfo["PlayerList"] = game.getLobbyList()
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
	gameInfo["PlayerList"] = target.getLobbyList()

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
	gameInfo["PlayerList"] = game.getTeamList()
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
		gameInfo["PlayerList"] = game.getPlayerUpdate()
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
	}
}