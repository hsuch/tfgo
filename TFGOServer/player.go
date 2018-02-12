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

	Status PlayerStatus
	StatusTimer *time.Timer

	Health int
	Armor int

	Weapon Weapon
	Inventory map[string]Pickup

	Location Location
}

func (p *Player) handleLoc(game *Game, loc Location) {

}

func (p *Player) fire(game *Game, dir Direction) {

}

func (p *Player) takeHit(wep Weapon) {

}