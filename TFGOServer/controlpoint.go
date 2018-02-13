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

	ControllingTeam Allegiance
	// number in [-7, 7] indicating capture/decapture progress.
	// hitting -7 or 7 from neutral ownership yields control
	// control point to red or blue, respectively. hitting 0
	// from team ownership neutralizes control point.
	CaptureProgress int
}

func (cp *ControlPoint) updateStatus(game *Game) {

}

func (cp *ControlPoint) inRange(loc Location) bool {
	if loc.getDistance(cp.Location) <= cp.Radius {
		return true
	}
	return false
}