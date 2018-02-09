package TFGOServer

import (
	"time"
	"math"
)

var SWORD = Weapon{
	Name: "Sword",
	Damage: 25,
	Spread: 2*math.Pi,
	Range: 10,
	ClipSize: 42,
	ShotReload: time.Second * 0,
	ClipReload: time.Second * 0,
}

type Weapon struct {
	Name string

	Damage int
	Spread float64
	Range float64

	ClipSize int
	ShotReload time.Duration
	ClipReload time.Duration
}
