package main

// clienthandler.go: primary goroutine serving each connected client

import (
	"net"
	"encoding/json"
)

// each incoming client message must conform to this format.
// for a detailed description of each message and its purpose,
// see https://github.com/hsuch/tfgo/wiki/Network-Messages
type ClientMessage struct {
	Action string
	Data map[string]interface{}
}

// the goroutine that serves a connected client. repeatedly
// processes messages and sends appropriate responses.
func serveClient(conn net.Conn) {
	defer conn.Close()

	// the current player and game
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
			g, p = createGame(conn, msg.Data)
			sendPlayerListUpdate(g)
		case "ShowGames":
			p = createPlayer(conn, msg.Data["Name"].(string), msg.Data["Icon"].(string))
			sendAvailableGames(p)
		case "ShowGameInfo":
			sendGameInfo(p, msg.Data["GameID"].(string))
		case "JoinGame":
			g = p.joinGame(msg.Data["GameID"].(string))
		case "StartGame":
			g.start()
		case "LocationUpdate":
			p.handleLoc(g, Location{X: 1, Y: 1})
		case "Fire":
			p.fire(g, weapons[msg.Data["Weapon"].(string)], msg.Data["Direction"].(float64))
		}
	}
}
