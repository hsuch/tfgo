package main

// fire.go: functions for handling firing weapon

import (
	"math"
)

// returns the dot product of two Direction vectors
func dot(v, w Direction) float64 {
	return v.X * w.X + v.Y * w.Y
}

// returns the magnitude of Direction vector v
func (v Direction) magnitude() float64 {
	return math.Sqrt(dot(v, v))
}

// determines whether the weapon fired from src in dir direction can hit a player at dst
// if it can, returns the distance from src to dst, otherwise, returns MaxFloat64
func (w Weapon) canHit(src, dst Location, dir Direction) float64 {
	target := Direction{dst.X - src.X, dst.Y - src.Y}
	dist := target.magnitude()
	angle := math.Acos(dot(target, dir) / (dist * dir.magnitude()))
	if angle <= w.Spread / 2 && dist <= w.Range {
		return dist
	} else {
		return math.MaxFloat64
	}
}

// fire the player's weapon at the given angle
func (p *Player) fire(game *Game, wep Weapon, angle float64) {
	if p.Status != NORMAL {
		return
	}
	//calculate direction vector of shot
	var dir Direction
	if angle == 0 || angle == 180 {
		dir = Direction{0,1}
	} else {
		dir.X = 1
		dir.Y = math.Cos(angle * math.Pi / 180) / math.Sin(angle * math.Pi / 180)
	}
	if angle >= 180 {
		dir.X *= -1
		dir.Y *= -1
	}
	minDist := math.MaxFloat64
	var closestP *Player
	// loop through list of enemies and find nearest hit
	for _, other := range game.Players {
		if p.Team != other.Team && other.Status != RESPAWNING {
			currDist := wep.canHit(p.Location, other.Location, dir)
			if currDist < minDist {
				closestP = other
				minDist = currDist
			}
		}
	}
	// if an enemy is hit, take the appropriate action
	if closestP != nil {
		closestP.takeHit(game, wep)
	}
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
	} else {
		sendTakeHit(p)
	}
}

func (p *Player) awaitRespawn(game *Game) {
	p.Status = RESPAWNING
	if p.OccupyingPoint != nil {
		if p.Team == game.RedTeam {
			p.OccupyingPoint.RedCount--
		} else if p.Team == game.BlueTeam {
			p.OccupyingPoint.BlueCount--
		}
		p.OccupyingPoint = nil
	}
	sendStatusUpdate(p, "Respawn")
	if !isTesting {
		<-p.StatusTimer.C
	}
	p.respawn(game)
}

func (p *Player) respawn(game *Game) {
	p.Status = NORMAL
	p.StatusTimer = nil
	p.Health = 100
	p.Armor = 0
	p.Inventory = nil
	if !inRange(p.Location, p.Team.Base, p.Team.BaseRadius) {
		go p.awaitRespawn(game)
	} else {
		sendStatusUpdate(p, "Respawned")
	}
}
