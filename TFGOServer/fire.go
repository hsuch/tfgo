package main

// fire.go: functions for handling firing weapon

import (
	"math"
	"time"
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
	dist := math.Max(target.magnitude(), 0.01) // account for exact same locations to prevent division by 0
	angle := math.Acos(dot(target, dir) / (dist * dir.magnitude()))
	if angle <= w.Spread / 2 && dist <= w.Range {
		return dist
	} else {
		return math.MaxFloat64
	}
}

// fire the player's weapon at the given angle
func (p *Player) fire(game *Game, wep Weapon, angle float64) {
	if game == nil || game.Status != PLAYING || p.Status != NORMAL {
		return
	}
	// calculate direction vector of shot
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
		p.Armor = p.Armor - wep.Damage
	} else {
		splash := wep.Damage - p.Armor
		p.Armor = 0
		p.Health = p.Health - splash
	}
	if p.Health <= 0 {
		p.Health = 0
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
	} else {
		sendVitalsUpdate(p)
	}
}

func (p *Player) awaitRespawn(game *Game) {
	p.StatusTimer = time.NewTimer(RESPAWNTIME())
	for {
		<- p.StatusTimer.C

		// account for early timer termination due to game end
		time.Sleep(10 * time.Millisecond)
		if p.StatusTimer == nil {
			break
		}

		if inRange(p.Location, p.Team.Base, p.Team.BaseRadius) {
			p.Status = NORMAL
			p.StatusTimer = nil
			p.Health = MAXHEALTH()
			p.Armor = 0
			sendStatusUpdate(p, "Respawned")
			break
		} else {
			p.StatusTimer = time.NewTimer(RESPAWNTIME())
		}
	}
}
