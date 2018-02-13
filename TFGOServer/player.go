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
