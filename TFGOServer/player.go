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

type Player struct {
	Name string
	Conn net.Conn
	Team *Team

	Status PlayerStatus
	StatusTimer *time.Timer

	Health int
	Armor int

	Weapon Weapon
	Inventory []Pickup

	Location Location
}

func (p *Player) handleLoc(loc Location) {

}

func (p *Player) fire(dir Direction) {

}

func (p *Player) takeHit(wep Weapon) {

}