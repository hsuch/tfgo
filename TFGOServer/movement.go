package main

// movement.go: functions for handling player movement

import (
	"time"
	"math"
)

/*
 * handleTimer - sets the out-of-bounds timer and kills the player if it fires
 *
 * p: Player whose life is at stake (is out of bounds)
 *
 * game: the Game struct containing all game-related information
 *
 * Returns: Nothing
 */
func handleOutOfBoundsTimer(game *Game, p *Player) {
	p.StatusTimer = time.AfterFunc(10 * time.Second, func() {
		p.awaitRespawn(game)
	})
}

/*
 * inBounds - check if location is within game boundaries
 *
 * game: the Game struct containing all game-related information
 *
 * loc: new Location to be checked
 *
 * Returns: True if loc is within game boundaries, false otherwise
 */
func inBounds(game *Game, loc Location) bool {
	intersections := 0.0
	for _, val := range game.Boundaries {
		t := val.P.Y + val.D.Y * (loc.X - val.P.X) / val.D.X - loc.Y
		s := (loc.X - val.P.X) / val.D.X
		if t >= 0 && s <= 1 {
			intersections++
		}
	}
	if math.Mod(intersections, 2) == 1 {
		return true
	} else {
		return false
	}
}

/*
 * handleLoc - sets a player's new location, checking for bounds and updating
 *             control points
 *
 * p: Player whose Location is being updated
 *
 * game: the Game struct containing all game-related information
 *
 * loc: new Location of Player p
 *
 * Returns: Nothing
 */
func (p *Player) handleLoc(game *Game, loc Location, orientation float64) {
	p.Location = loc
	p.Orientation = orientation

	// Check if player is in/out of bounds and handle accordingly
	if !inBounds(game, p.Location) && p.Status == NORMAL {
		p.Status = OUTOFBOUNDS
		go handleOutOfBoundsTimer(game, p)
	} else if inBounds(game, p.Location) && p.Status == OUTOFBOUNDS {
		p.Status = NORMAL
		p.StatusTimer.Stop()
	}

	// Get player's new control point (player status MUST be NORMAL)
	var newCP *ControlPoint = nil
	if p.Status == NORMAL {
		for _, cp := range game.ControlPoints {
			if cp.inRange(p.Location) {
				newCP = cp
				break
			}
		}
	}

	// Set new control point and change team counts at control points accordingly
	if newCP != nil && newCP != p.OccupyingPoint {
		if p.Team == game.RedTeam {
			newCP.RedCount++
		} else if p.Team == game.BlueTeam {
			newCP.BlueCount++
		}
	}
	if p.OccupyingPoint != nil && newCP != p.OccupyingPoint{
		if p.Team == game.RedTeam {
			p.OccupyingPoint.RedCount--
		} else if p.Team == game.BlueTeam {
			p.OccupyingPoint.BlueCount--
		}
	}
	p.OccupyingPoint = newCP
}

/*
 * updateStatus - updates a capture points's capture progress and controlling
 *                team, as well as and team points if appropriate
 *
 * cp: ControlPoint to be updated
 *
 * game: the Game struct containing all game-related information
 *
 * Returns: Nothing
 */
func (cp *ControlPoint) updateStatus(game *Game) {
	// Update capture progress
	oldCaptureProgress := cp.CaptureProgress
	cp.CaptureProgress += cp.BlueCount - cp.RedCount

	// Change ControllingTeam to NEUTRAL if either
	// (1) CaptureProgress hits 0 OR
	// (2) sign of CaptureProgress changes
	if cp.CaptureProgress == 0 || oldCaptureProgress ^ cp.CaptureProgress < 0  {
		cp.ControllingTeam = nil
	}

	// Set appropriate ControllingTeam
	if cp.CaptureProgress >= 7 {
		cp.ControllingTeam = game.BlueTeam
		cp.CaptureProgress = 7
	} else if cp.CaptureProgress <= -7 {
		cp.ControllingTeam = game.RedTeam
		cp.CaptureProgress = -7
	}

	// Add a point to the team controlling this control point
	cp.ControllingTeam.Points++
}

func (l1 Location) getDistance(l2 Location) float64 {
	first := math.Pow(float64(l2.X-l1.X), 2)
	second := math.Pow(float64(l2.Y-l1.Y), 2)
	return math.Sqrt(first + second)
}

func (cp *ControlPoint) inRange(loc Location) bool {
	if loc.getDistance(cp.Location) <= cp.Radius {
		return true
	}
	return false
}
