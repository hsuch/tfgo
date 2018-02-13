package main

import (
	"net"
	"time"
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

const initStatus = NORMAL
const initHealth = 100
const initArmor = 0
const initPickup = nil

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

	Status PlayerStatus
	StatusTimer *time.Timer

	Health int
	Armor int

	Inventory map[string]Pickup

	Location Location
}

// returns true if loc is within the boundaries of the game
func inBounds(game *Game, loc Location) bool {
	intersections := 0
	_, val := range game.Boundaries {
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

func (p *Player) handleLoc(game *Game, loc Location) {

}

// decides whether the shot hits anyone, and if so, calls takeHit()
func (p *Player) fire(game *Game, wep Weapon, dir Direction) {
	min_dist := math.MaxFloat64
	var closest_p *Player
	var enemies []Player
	if p.Allegiance == RED {
		enemies = game.BlueTeam.GetPlayerLocs()
	} else if p.Allegiance == BLUE {
		enemies = game.RedTeam.GetPlayerLocs()
	} else {
		enemies = nil
	}

	//loop through list of enemies and find nearest hit
	_, val := range enemies {
		curr_dist := wep.canHit(p.Location, val.Location, dir)
		if curr_dist < min_dist {
			closest_p = val
			min_dist = curr_dist
		}
	}

	//if an enemy is hit, take the appropriate action
	if closest_p != nil {
		closest_p.takeHit(wep)
	}
	//will we want to communicate to the client that they've hit someone?
}

/* Input: Weapon dealing damage
   Output: Void, runs goroutine to
   place player in respawn queue if killed.
*/
func (p *Player) takeHit(wep Weapon) bool {
	if wep.Damage <= p.armor {
		p.armor -= wep.Damage
	}
	else {
		splash := wep.Damage - p.armor
		p.health -= splash
	}

	if p.health <= 0 {
		go awaitRespawn ()
	}
}

/* Called by awaitRespawn () goroutine after timer reset */
func (p *player) respawn () {
	p.status = RESPAWNING
	p.health = initHealth
	p.armor = initArmor
	p.pickup = initPickup
	p.location = p.Team.Base
}
