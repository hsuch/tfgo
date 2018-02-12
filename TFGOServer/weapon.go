package main

import (
	"time"
	"math"
)

type Direction struct {
	X float64
	Y float64
}

// each of the available weapons is defined as a globally
// accessible variable
var SWORD = Weapon {
	Name: "Sword",
	Damage: 25,
	Spread: 2*math.Pi,
	Range: 1,
	ClipSize: 500,
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

func (w Weapon) canHit(src, dst Location, dir Direction) bool {
	return false
}