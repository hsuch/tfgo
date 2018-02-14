package main

import (
	"time"
	"math"
)

type Direction struct {
	X float64
	Y float64
}

//returns the dot product of two Direction vectors
func dot(v, w Direction) float64 {
	return v.X * w.X + v.Y * w.Y
}

//returns the magnitude of Direction vector v
func (v Direction) magnitude() float64 {
	return math.Sqrt(dot(v, v))
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

var SHOTGUN = Weapon {
	Name: "Shotgun",
	Damage: 25,
	Spread: math.Pi/2,
	Range: 3,
	ClipSize: 2,
	ShotReload: time.Millisecond * 500,
	ClipReload: time.Second * 3,
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

// determines whether the weapon fired from src in dir direction can hit a player at dst
// if it can, returns the distance from src to dst, if not returns MaxFloat64
func (w Weapon) canHit(src, dst Location, dir Direction) float64 {
	target := Direction{dst.X - src.X, dst.Y - src.Y}
	dist := target.magnitude()
	angle := math.Acos(dot(target, dir) / (dist * dir.magnitude()))
	if angle <= w.Spread && dist <= w.Range {
		return dist
	} else {
		return math.MaxFloat64
	}
}
