package main

type ControlPoint struct {
	ID string

	Location Location
	Radius float64

	PayloadPath [2]Location // start, end
	PayloadLoc Location

	// number of currently occupying players from each team
	RedCount int
	BlueCount int

	// number in [-7, 7] indicating capture/decapture progress.
	// hitting -7 or 7 from neutral ownership yields control
	// control point to red or blue, respectively. hitting 0
	// from team ownership neutralizes control point.
	CaptureProgress int
	ControllingTeam Allegiance
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
		cp.ControllingTeam = NEUTRAL
	}

	// Set appropriate ControllingTeam
	if cp.CaptureProgress >= 7 {
		cp.ControllingTeam = BLUE
		cp.CaptureProgress = 7
	} else if cp.CaptureProgress <= -7 {
		cp.ControllingTeam = RED
		cp.CaptureProgress = -7
	}

	// Add a point to the team controlling this control point
	if cp.ControllingTeam == RED {
		game.RedTeam.Points++
	} else if cp.ControllingTeam == BLUE {
		game.BlueTeam.Points++
	}
}

func (cp *ControlPoint) inRange(loc Location) bool {
	if loc.getDistance(cp.Location) <= cp.Radius {
		return true
	}
	return false
}