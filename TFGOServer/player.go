package main

import (
	"net"
	"time"
)

type PlayerStatus int
const (
	NORMAL PlayerStatus = iota
	OUTOFBOUNDS
	RESPAWNING
)

var PlayerStatusMap = map[PlayerStatus]string{
	NORMAL : "NORMAL",
	OUTOFBOUNDS : "OUTOFBOUNDS",
	RESPAWNING : "RESPAWNING",
}

type Allegiance int
const (
	RED Allegiance = iota
	NEUTRAL
	BLUE
)

type Player struct {
	Name string
	Conn net.Conn
	Team Allegiance
	ControlPoint *ControlPoint

	Status PlayerStatus
	StatusTimer *time.Timer

	Health int
	Armor int

	Weapon Weapon
	Inventory map[string]Pickup

	Location Location
}

/*
 * handleTimer - sets the out-of-bounds timer and kills the player if it fires
 *
 * p: player whose life is at stake (is out of bounds)
 *
 * game: the game struct containing all game-related information
 */
func handleOutOfBoundsTimer(game *Game, p *Player) {
	p.StatusTimer = time.AfterFunc(10 * time.Second, func() {
		p.awaitRespawn(game)
	})
}

/*
 * handleLoc - sets a player's new location, checking for bounds and updating
 *             control points
 *
 * p: player whose location is being updated
 *
 * game: the game struct containing all game-related information
 *
 * loc: new location of player p
 */
func (p *Player) handleLoc(game *Game, loc Location) {
	p.Location = loc

	// Check if player is in/out of bounds and handle accordingly
	if !inBounds(game, p.Location) && p.Status == NORMAL {
		p.Status = OUTOFBOUNDS
		go handleOutOfBoundsTimer(game, p)
	} else if inBounds(game, p.Location) && p.Status == OUTOFBOUNDS {
		p.Status = NORMAL
		p.StatusTimer.Stop()
	}

	// Get player's new control point (player status MUST be NORMAL)
	var newCP *ControlPoint = nil
	if p.Status == NORMAL {
		for _, cp := range game.ControlPoints {
			if cp.inRange(p.Location) {
				newCP = cp
				break
			}
		}
	}

	// Set new control point and change team counts at control points accordingly
	if newCP != nil && newCP != p.ControlPoint {
		if p.Team == RED {
			newCP.RedCount++
		} else if p.Team == BLUE {
			newCP.BlueCount++
		}
	}
	if p.ControlPoint != nil && newCP != p.ControlPoint {
		if p.Team == RED {
			p.ControlPoint.RedCount--
		} else if p.Team == BLUE {
			p.ControlPoint.BlueCount--
		}
	}
	p.ControlPoint = newCP
}

func (p *Player) fire(game *Game, dir Direction) {

}

func (p *Player) takeHit(wep Weapon) {

}