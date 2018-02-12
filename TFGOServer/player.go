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

type Allegiance int
const (
	RED Allegiance = iota
	NEUTRAL
	BLUE
)

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

func playerStatusString(ps PlayerStatus) string {
	switch ps {
	case NORMAL:
		return "NORMAL"
	case OUTOFBOUNDS:
		return "OUTOFBOUNDS"
	case RESPAWNING:
		return "RESPAWNING"
	default:
		return ""
	}
}

func (p *Player) handleLoc(game *Game, loc Location) {

}

func (p *Player) fire(game *Game, wep Weapon, dir Direction) {

}

func (p *Player) takeHit(wep Weapon) {

}
