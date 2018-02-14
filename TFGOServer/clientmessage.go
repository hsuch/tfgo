package main

// clientmessage.go: functions for building and sending messages to clients

func (p *Player) sender() {
	for {
		msg := <-p.Chan
		p.Encoder.Encode(msg)
	}
}

func getLobbyList(game *Game) []map[string]string {
	var playerList []map[string]string
	for _, player := range game.RedTeam.Players {
		playerInfo := make(map[string]string)
		playerInfo["Name"] = player.Name
		playerInfo["Icon"] = player.Icon
		playerList = append(playerList, playerInfo)
	}

	return playerList
}

func sendPlayerListUpdate(game *Game) {
	for _, player := range game.RedTeam.Players {
		msg := map[string]interface{} {
			"Type" : "PlayerListUpdate",
			"Data" : getLobbyList(game),
		}
		player.Chan <- msg
	}
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
			gameInfo["PlayerList"] = getLobbyList(game)
			gameList = append(gameList, gameInfo)
		}
	}

	msg := map[string]interface{} {
		"Type" : "AvailableGames",
		"Data" : gameList,
	}
	player.Chan <- msg
}