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

func (p *Player) handleLoc(game *Game, loc Location) {

}

func (p *Player) fire(game *Game, wep Weapon, dir Direction) {

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