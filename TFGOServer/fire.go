package main

import "math"

//returns the dot product of two Direction vectors
func dot(v, w Direction) float64 {
	return v.X * w.X + v.Y * w.Y
}

//returns the magnitude of Direction vector v
func (v Direction) magnitude() float64 {
	return math.Sqrt(dot(v, v))
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

/*
 * fire() - determines whether the shot hits anyone, if so calls takeHit()
 *
 * p: Player who fired the shot
 *
 * game: the Game struct containing all game-related information
 *
 * wep: the Weapon used to fire the shot
 *
 * dir: Direction vector of the shot
 *
 * Returns: Nothing
 */
func (p *Player) fire(game *Game, wep Weapon, dir Direction) {
	minDist := math.MaxFloat64
	var closestP *Player
	var enemies map[string]*Player
	if p.Team == RED {
		enemies = game.BlueTeam.Players
	} else if p.Team == BLUE {
		enemies = game.RedTeam.Players
	} else {
		enemies = nil
	}

	// loop through list of enemies and find nearest hit
	for _, enemy := range enemies {
		if enemy.Status == NORMAL {
			currDist := wep.canHit(p.Location, enemy.Location, dir)
			if currDist < minDist {
				closestP = enemy
				minDist = currDist
			}
		}
	}

	// if an enemy is hit, take the appropriate action
	if closestP != nil {
		closestP.takeHit(game, wep)
	}
	// will we want to communicate to the client that they've hit someone?
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
	p.StatusTimer = nil
	p.Health = 100
	p.Armor = 0
	p.Inventory = nil
	// if player not in base...
}
