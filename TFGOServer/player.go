package main

import (
	"net"
	"time"
	"math"
)

type Allegiance int
const (
	RED Allegiance = iota
	NEUTRAL
	BLUE
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

type Player struct {
	Name string
	Icon string
	Conn net.Conn
	Team Allegiance
	ControlPoint *ControlPoint

	Status PlayerStatus
	StatusTimer *time.Timer

	Health int
	Armor int

	Inventory map[string]Pickup

	Location Location
}

/*
 * handleTimer - sets the out-of-bounds timer and kills the player if it fires
 *
 * p: Player whose life is at stake (is out of bounds)
 *
 * game: the Game struct containing all game-related information
 *
 * Returns: Nothing
 */
func handleOutOfBoundsTimer(game *Game, p *Player) {
	p.StatusTimer = time.AfterFunc(10 * time.Second, func() {
		p.awaitRespawn(game)
	})
}

/*
 * inBounds - check if location is within game boundaries
 *
 * game: the Game struct containing all game-related information
 *
 * loc: new Location to be checked
 *
 * Returns: True if loc is within game boundaries, false otherwise
 */
func inBounds(game *Game, loc Location) bool {
	intersections := 0.0
	for _, val := range game.Boundaries {
		t := (val.P.Y + val.D.Y * (loc.X - val.P.X) / val.D.X - loc.Y)
		s := (loc.X - val.P.X) / val.D.X
		if t >= 0 && s <= val.T {
			intersections++
		}
	}
	if math.Mod(intersections, 2) == 1 {
		return true
	} else {
		return false
	}
}

/*
 * handleLoc - sets a player's new location, checking for bounds and updating
 *             control points
 *
 * p: Player whose Location is being updated
 *
 * game: the Game struct containing all game-related information
 *
 * loc: new Location of Player p
 *
 * Returns: Nothing
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

/*
 * fire() - determines whether the shot hits anyone, if so calls takeHit()
 *
 * p: Player who fired the shot
 *
 * game: the Game struct containing all game-related information
 *
 * wep: the Weapon used to fire the shot
 *
 * dir: Direction vector of the shot
 *
 * Returns: Nothing
 */
func (p *Player) fire(game *Game, wep Weapon, dir Direction) {
	min_dist := math.MaxFloat64
	var closest_p *Player
	var enemies map[string]*Player
	if p.Team == RED {
		enemies = game.BlueTeam.Players
	} else if p.Team == BLUE {
		enemies = game.RedTeam.Players
	} else {
		enemies = nil
	}

	// loop through list of enemies and find nearest hit
	for _, enemy := range enemies {
		if enemy.Status == NORMAL {
			curr_dist := wep.canHit(p.Location, enemy.Location, dir)
			if curr_dist < min_dist {
				closest_p = enemy
				min_dist = curr_dist
			}
		}
	}

	// if an enemy is hit, take the appropriate action
	if closest_p != nil {
		closest_p.takeHit(game, wep)
	}
	// will we want to communicate to the client that they've hit someone?
}

func (p *Player) takeHit(game *Game, wep Weapon) {
	if wep.Damage <= p.Armor {
		p.Armor -= wep.Damage
	} else {
		splash := wep.Damage - p.Armor
		p.Armor = 0
		p.Health -= splash
	}

	if p.Health <= 0 {
		go p.awaitRespawn(game)
	}

	// send message to player informing hit
}

func (p *Player) awaitRespawn(game *Game) {
	p.Status = RESPAWNING
	<- p.StatusTimer.C
	p.respawn(game)
	// send message to player informing respawn
}

func (p *Player) respawn(game *Game) {
	p.Status = NORMAL
	p.Health = 100
	p.Armor = 0
	p.Inventory = nil
	// if player not in base...
}
