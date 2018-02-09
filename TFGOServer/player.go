package TFGOServer

import (
	"net"
	"time"
)

type Player struct {
	Name string
	Conn net.Conn
	Team Team

	Status int
	StatusTimer time.Timer

	Health int
	Armor int

	Weapon Weapon
	Inventory []Pickup

	Location Location
}
