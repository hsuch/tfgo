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

func (p *Player) handleLoc(game *Game, loc Location) {

}

func (p *Player) fire(game *Game, wep Weapon, dir Direction) {

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