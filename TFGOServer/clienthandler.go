package main

// clienthandler.go: primary goroutine serving each connected client

import (
	"net"
	"encoding/json"
	"fmt"
	"os"
)

var verbosity = 0

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

	fmt.Println("New connection received.")

	// the current player and game
	var p *Player
	var g *Game

	decoder := json.NewDecoder(conn)
	for {
		var msg ClientMessage
		if err := decoder.Decode(&msg); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read message: %s\nClosing connection.\n", err.Error())

			if p != nil {
				fmt.Printf("Player %s disconnected.\n", p.Name)
				p.leaveGame(g)
				close(p.Chan)
			}
			break
		}

		if verbosity > 1 || verbosity > 0 && msg.Action != "LocationUpdate" {
			rawJSON, _ := json.Marshal(msg)
			source := "client"
			if p != nil {
				source = p.Name
			}
			fmt.Printf("Received from %s:\n%s\n", source, prettyPrintJSON(rawJSON))
		}

		switch msg.Action {
		case "RegisterPlayer":
			p = createPlayer(conn, msg.Data["Name"].(string), msg.Data["Icon"].(string))
			sendPlayerID(p)
		case "ChangePlayerInfo":
			p.Name = msg.Data["Name"].(string)
			p.Icon = msg.Data["Icon"].(string)
		case "CreateGame":
			g = p.createGame(conn, msg.Data)
			sendPlayerListUpdate(g)
		case "ShowGames":
			sendAvailableGames(p)
		case "ShowGameInfo":
			sendGameInfo(p, msg.Data["GameID"].(string))
		case "JoinGame":
			g = p.joinGame(msg.Data["GameID"].(string), msg.Data["Password"].(string))
		case "LeaveGame":
			p.leaveGame(g)
		case "StartGame":
			g.start()
		case "LocationUpdate":
			loc := msg.Data["Location"].(map[string]interface{})
			p.updateLocation(g, Location{X: loc["X"].(float64), Y: loc["Y"].(float64)}, msg.Data["Orientation"].(float64))
		case "Fire":
			p.fire(g, weapons[msg.Data["Weapon"].(string)], msg.Data["Direction"].(float64))
		}
	}
}
