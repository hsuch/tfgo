package main

import (
	"net"
	"encoding/json"
)

type ClientMessage struct {
	Action string
	Data map[string]interface{}
}

func serveClient(conn net.Conn) {
	defer conn.Close()

	var p *Player
	var g *Game

	decoder := json.NewDecoder(conn)
	for {
		var msg ClientMessage
		if err := decoder.Decode(&msg); err != nil {
			break
		}

		switch msg.Action {
		case "CreateGame":
			g, p = createNewGame(conn, msg.Data)
		case "ShowGames":
			sendAvailableGames(p)
		case "ShowGameInfo":
			sendGameInfo(msg.Data["GameID"].(string))
		case "JoinGame":
		case "StartGame":
		case "LocationUpdate":
		case "Fire":
			p.fire(g, weapons[msg.Data["Weapon"].(string)], msg.Data["Direction"].(float64))
		}
	}
}
