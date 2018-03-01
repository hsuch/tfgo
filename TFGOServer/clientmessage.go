package main

import (
	"time"
	"reflect"
	"bytes"
	"encoding/json"
	"fmt"
)

// clientmessage.go: functions for building and sending messages to clients

// converts a ClientMessage into formatted JSON that can be fed to printing functions
func prettyPrintJSON(rawJSON []byte) string {
	var out bytes.Buffer
	json.Indent(&out, rawJSON, "", "    ")
	return string(out.Bytes())
}

// goroutine responsible for delivering messages to players
func (p *Player) sender() {
	for {
		msg, open := <-p.Chan
		if open {
			if verbosity > 1 || verbosity > 0 && msg["Type"] != "GameUpdate" {
				rawJSON, _ := json.Marshal(msg)
				fmt.Printf("Sending to %s:\n%s\n", p.Name, prettyPrintJSON(rawJSON))
			}
			if !isTesting {
				p.Encoder.Encode(msg)
			}
		} else {
			return
		}
	}
}

// prevents sending to closed channels, as we close and send in different places
func (p *Player) safeSend(msg map[string]interface{}) {
	defer func() { recover() }()
	p.Chan <- msg
}

// deliver a message to all players
func (g *Game) broadcast(msg map[string]interface{}) {
	for _, player := range g.Players {
		player.safeSend(msg)
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
			} else if field == "Location" {
				playerInfo[field] = player.Location.locationToDegrees()
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
		vertex["X"] = meterToDegree(boundary.P.X)
		vertex["Y"] = meterToDegree(boundary.P.Y)
		vertices = append(vertices, vertex)
	}

	return vertices
}

func (t *Team) getLocInfo() map[string]interface{} {
	return map[string]interface{} {
		"Location" : t.Base.locationToDegrees(),
		"Radius" : t.BaseRadius, //in meters because the app team wants it that way
	}
}

func (cp *ControlPoint) getLocInfo() map[string]interface{} {
	return map[string]interface{} {
		"ID" : cp.ID,
		"Location" : cp.Location.locationToDegrees(),
		"Radius" : cp.Radius, //also in meters
	}
}

func (p *PickupSpot) getInfo() map[string]interface{} {
	var pickupType string
	var pickupAmt int
	if ap, ok := p.Pickup.(ArmorPickup); ok {
		pickupType = "Armor"
		pickupAmt = ap.AP
	} else if hp, ok := p.Pickup.(HealthPickup); ok {
		pickupType = "Health"
		pickupAmt = hp.HP
	} else if wp, ok := p.Pickup.(WeaponPickup); ok {
		pickupType = wp.WP.Name
		pickupAmt = 1
	}

	return map[string]interface{} {
		"Location" : p.Location,
		"Type" : pickupType,
		"Amount" : pickupAmt,
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
		cpInfo["ID"] = cp.ID
		cpInfo["Location"] = cp.Location.locationToDegrees()
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
		if game.Status == CREATING && len(game.Players) < game.PlayerLimit {
			gameInfo := make(map[string]interface{})
			gameInfo["ID"] = game.HostID
			gameInfo["Name"] = game.Name
			if game.Password == "" {
				gameInfo["HasPassword"] = false
			} else {
				gameInfo["HasPassword"] = true
			}
			gameInfo["Mode"] = modeToString[game.Mode]
			gameInfo["Location"] = game.findCenter().locationToDegrees()
			gameInfo["PlayerList"] = game.getPlayerInfo([]string{"Name", "Icon"})
			gameList = append(gameList, gameInfo)
		}
	}

	msg := map[string]interface{} {
		"Type" : "AvailableGames",
		"Data" : gameList,
	}
	player.safeSend(msg)
}

func sendGameInfo(player *Player, gameID string) {
	target := games[gameID]
	gameInfo := make(map[string]interface{})
	gameInfo["ID"] = target.HostID
	gameInfo["Description"] = target.Description
	gameInfo["PlayerLimit"] = target.PlayerLimit
	gameInfo["PointLimit"] = target.PointLimit
	gameInfo["TimeLimit"] = target.TimeLimit.Minutes()
	gameInfo["Boundaries"] = target.getBoundaryVertices()
	gameInfo["PlayerList"] = target.getPlayerInfo([]string{"Name", "Icon"})

	msg := map[string]interface{} {
		"Type" : "GameInfo",
		"Data" : gameInfo,
	}
	player.safeSend(msg)
}

func sendJoinGameError(player *Player, error string) {
	msg := map[string]interface{} {
		"Type" : "JoinGameError",
		"Data" : error,
	}
	player.safeSend(msg)
}

func sendLeaveGame(game *Game) {
	msg := map[string]interface{}{
		"Type" : "LeaveGame",
	}

	for id, player := range game.Players {
		if id != game.HostID {
			player.safeSend(msg)
		}
	}
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
	var pickupInfo []map[string]interface{}
	for _, pickup := range game.Pickups {
		pickupInfo = append(pickupInfo, pickup.getInfo())
	}
	gameInfo["Pickups"] = pickupInfo
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

func sendStatusUpdate(player *Player, status string) {
	msg := map[string]interface{} {
		"Type" : "StatusUpdate",
		"Data" : status,
	}
	player.safeSend(msg)
}

func sendVitalsUpdate(player *Player) {
	msg := map[string]interface{} {
		"Type" : "VitalsUpdate",
		"Data" : map[string]int {
			"Health" : player.Health,
			"Armor" : player.Armor,
		},
	}
	player.safeSend(msg)
}

func sendPickupUpdate(pickup *PickupSpot, game *Game) {
	msg := map[string]interface{} {
		"Type" : "PickupUpdate",
		"Data" : map[string]interface{} {
			"Location" : pickup.Location,
			"Available" : pickup.Available,
		},
	}
	game.broadcast(msg)
}

func sendAcquireWeapon(player *Player, weapon Weapon) {
	msg := map[string]interface{} {
		"Type" : "AcquireWeapon",
		"Data" : weaponToString[weapon],
	}
	player.safeSend(msg)
}

func sendGameOver(game *Game) {
	winner := "None"
	if game.Mode == PAYLOAD {
		payloadLoc := game.ControlPoints["Payload"].Location
		if inRange(payloadLoc, game.RedTeam.Base, game.RedTeam.BaseRadius) {
			winner = game.BlueTeam.Name
		} else if inRange(payloadLoc, game.BlueTeam.Base, game.BlueTeam.BaseRadius) {
			winner = game.RedTeam.Name
		}
	} else {
		if game.RedTeam.Points > game.BlueTeam.Points {
			winner = game.RedTeam.Name
		} else if game.RedTeam.Points < game.BlueTeam.Points {
			winner = game.BlueTeam.Name
		}
	}

	msg := map[string]interface{} {
		"Type" : "GameOver",
		"Data" : winner,
	}
	game.broadcast(msg)
}
