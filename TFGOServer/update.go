package main

// update.go: functions for handling updates

import (
	"time"
	"math"
)

// returns distance between l1 and l2
func distance(l1, l2 Location) float64 {
	first := math.Pow(float64(l2.X - l1.X), 2)
	second := math.Pow(float64(l2.Y - l1.Y), 2)
	return math.Sqrt(first + second)
}

// checks whether l1 is within dist distance of l2
func inRange(l1, l2 Location, dist float64) bool {
	return distance(l1, l2) <= dist
}

// check if loc is within game boundaries
func inGameBounds(game *Game, loc Location) bool {
	intersections := 0.0
	for _, val := range game.Boundaries {
		if loc.X == val.P.X && loc.Y == val.P.Y {
			return true
		}
		t := val.P.Y + (val.D.Y * (loc.X - val.P.X) / val.D.X) - loc.Y
		s := (loc.X - val.P.X) / val.D.X
		if t > 0.00000000001 && s >= 0 && s <= 1 {
			intersections += 1
		}
	}
	return math.Mod(intersections, 2.0) == 1.0
}

// in PAYLOAD mode, moves the payload towards the base of the controlling team
func movePayload(game *Game) {
	cp := game.ControlPoints["Payload"]
	if cp.ControllingTeam == nil {
		return
	}
	var t float64
	if cp.ControllingTeam == game.BlueTeam {
		t = -1
	} else {
		t = 1
	}

	cp.Location.X += t * game.PayloadPath.X * game.PayloadSpeed
	cp.Location.Y += t * game.PayloadPath.Y * game.PayloadSpeed
	return
}

 // set player location, updating respawn, out-of-bounds, and control point
 // info as necessary
func (p *Player) updateLocation(game *Game, loc Location, orientation float64) {
	p.Location = Location{X: degreeToMeter(loc.X), Y: degreeToMeter(loc.Y)}
	p.Orientation = orientation

	// no information should be updated if the game has not yet started
	if game.Status == CREATING {
		return
	}

	// start respawn timer if they just returned to their base after dying
	if p.Status == RESPAWNING && p.StatusTimer == nil && inRange(p.Location, p.Team.Base, p.Team.BaseRadius) {
		p.awaitRespawn(game)
	}

	// handle entering/exiting game boundaries
	if p.Status == NORMAL && !inGameBounds(game, p.Location) {
		p.Status = OUTOFBOUNDS
		p.StatusTimer = time.AfterFunc(OUTOFBOUNDSTIME(), func() {
			p.Status = RESPAWNING
			sendStatusUpdate(p, "Respawn")
		})
		sendStatusUpdate(p, "OutOfBounds")
	} else if p.Status == OUTOFBOUNDS && inGameBounds(game, p.Location) {
		p.Status = NORMAL
		p.StatusTimer.Stop()
		p.StatusTimer = nil
		sendStatusUpdate(p, "BackInBounds")
	}

	// check if player if on a pickup
	if p.Status == NORMAL {
		for _, pickup := range game.Pickups {
			if inRange(p.Location, pickup.Location, PICKUPRADIUS()) {
				pickup.consumePickup(p, game)
				break
			}
		}
	}

	// get player's new control point (player status MUST be NORMAL)
	var newCP *ControlPoint = nil
	if p.Status == NORMAL {
		for _, cp := range game.ControlPoints {
			if inRange(p.Location, cp.Location, cp.Radius) {
				newCP = cp
				break
			}
		}
	}

	// set new control point and change team counts at control points accordingly
	if newCP != nil && newCP != p.OccupyingPoint {
		if p.Team == game.RedTeam {
			newCP.RedCount++
		} else if p.Team == game.BlueTeam {
			newCP.BlueCount++
		}
	}
	if p.OccupyingPoint != nil && newCP != p.OccupyingPoint {
		if p.Team == game.RedTeam {
			p.OccupyingPoint.RedCount--
		} else if p.Team == game.BlueTeam {
			p.OccupyingPoint.BlueCount--
		}
	}
	p.OccupyingPoint = newCP
}

// updates control point capture progress, controlling team, and team points
func (cp *ControlPoint) updateStatus(game *Game) {
	// update capture progress
	oldCaptureProgress := cp.CaptureProgress
	delta := cp.BlueCount - cp.RedCount
	if delta > 3 {
		delta = 3
	} else if delta < -3 {
		delta = -3
	}
	cp.CaptureProgress += delta

	// remove ControllingTeam if either
	// (1) CaptureProgress hits 0 OR
	// (2) sign of CaptureProgress changes
	if cp.CaptureProgress == 0 || oldCaptureProgress * cp.CaptureProgress < 0  {
		cp.ControllingTeam = nil
	}

	// set appropriate ControllingTeam
	if cp.CaptureProgress >= 7 {
		cp.ControllingTeam = game.BlueTeam
		cp.CaptureProgress = 7
	} else if cp.CaptureProgress <= -7 {
		cp.ControllingTeam = game.RedTeam
		cp.CaptureProgress = -7
	}

	// modify data for the controlling team
	if cp.ControllingTeam != nil {
		if game.Mode == PAYLOAD {
			movePayload(game)
			payloadLoc := game.ControlPoints["Payload"].Location
			if inRange(payloadLoc, game.RedTeam.Base, game.RedTeam.BaseRadius) ||
				inRange(payloadLoc, game.BlueTeam.Base, game.BlueTeam.BaseRadius) {
				game.Timer.Stop()
				sendGameUpdate(game)
				game.stop()
			}
		} else {
			cp.ControllingTeam.Points++
			if cp.ControllingTeam.Points == game.PointLimit {
				game.Timer.Stop()
				sendGameUpdate(game)
				game.stop()
			}
		}
	}
}

// repeatedly calls updateStatus
func (cp *ControlPoint) updateTicker(game *Game) {
	for game.Status == PLAYING {
		cp.updateStatus(game)
		time.Sleep(time.Second)
	}
}
