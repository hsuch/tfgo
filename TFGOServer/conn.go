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

	d := json.NewDecoder(conn)
	for {
		var msg ClientMessage
		if err := d.Decode(&msg); err != nil {
			break
		}

		switch msg.Action {
		case "CreateGame":
			g, p = createNewGame(msg.Data)
		case "ShowGames":
			sendAvailableGames(p)
		case "ShowGameInfo":
		case "JoinGame":
		case "StartGame":
		case "LocationUpdate":
		case "Fire":
			p.fire(g, weapons[msg.Data["Weapon"].(string)], msg.Data["Direction"].(float64))
		}
	}
}
