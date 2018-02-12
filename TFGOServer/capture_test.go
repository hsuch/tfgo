package main

import "testing"

func TestHandleLoc(t *testing.T) {
	g := makeSampleGame()
	cp := g.ControlPoints["CP2"]
	oliver := getOliver(g)

	// player moves out of bounds
	oliver.handleLoc(g, Location{420, 170})
	if oliver.Status != OUTOFBOUNDS {
		t.Errorf("TestHandleLoc(1) failed, expected Status OUTOFBOUNDS, got Status %s", playerStatusString(oliver.Status))
	}

	// player exits control point
	expBlueCount := cp.BlueCount - 1
	oliver.handleLoc(g, Location{0, 0})
	if cp.BlueCount != expBlueCount {
		t.Errorf("TestHandleLoc(3) failed, expected BlueCount %d, got BlueCount %d", expBlueCount, cp.BlueCount)
	}

	// player enters control point
	expBlueCount = cp.BlueCount + 1
	oliver.handleLoc(g, cp.Location)
	if cp.BlueCount != expBlueCount {
		t.Errorf("TestHandleLoc(2) failed, expected BlueCount %d, got BlueCount %d", expBlueCount, cp.BlueCount)
	}
}

func TestUpdateStatus(t *testing.T) {
	g := makeSampleGame()
	cp1 := g.ControlPoints["CP1"]
	cp2 := g.ControlPoints["CP2"]
	redTeam := g.RedTeam
	blueTeam := g.BlueTeam
	oliver := getOliver(g)

	// equal players from each team; check capture progress
	expCaptureProg := cp2.CaptureProgress
	oliver.handleLoc(g, Location{0, 0})
	cp2.updateStatus(g)
	if cp2.CaptureProgress != expCaptureProg {
		t.Errorf("TestUpdateStatus(2) failed, expected CaptureProgress %d, got CaptureProgress %d", expCaptureProg, cp2.CaptureProgress)
	}

	// blue players exceed red by one; check capture progress
	expCaptureProg = cp2.CaptureProgress + 1
	oliver.handleLoc(g, cp2.Location)
	cp2.updateStatus(g)
	if cp2.CaptureProgress != expCaptureProg {
		t.Errorf("TestUpdateStatus(1) failed, expected CaptureProgress %d, got CaptureProgress %d", expCaptureProg, cp2.CaptureProgress)
	}

	// both control points uncontrolled; check team points
	cp1.ControllingTeam = NEUTRAL
	cp2.ControllingTeam = NEUTRAL
	expRedPoints := redTeam.Points
	expBluePoints := blueTeam.Points
	cp1.updateStatus(g)
	cp2.updateStatus(g)
	if redTeam.Points != expRedPoints {
		t.Errorf("TestUpdateStatus(3) failed, expected Red Team Points %d, got Red Team Points %d", expRedPoints, redTeam.Points)
	}
	if blueTeam.Points != expBluePoints {
		t.Errorf("TestUpdateStatus(3) failed, expected Blue Team Points %d, got Blue Team Points %d", expBluePoints, blueTeam.Points)
	}

	// both control points controlled by blue; check team points
	cp1.ControllingTeam = BLUE
	cp2.ControllingTeam = BLUE
	expRedPoints = redTeam.Points
	expBluePoints = blueTeam.Points + 2
	cp1.updateStatus(g)
	cp2.updateStatus(g)
	if redTeam.Points != expRedPoints {
		t.Errorf("TestUpdateStatus(4) failed, expected Red Team Points %d, got Red Team Points %d", expRedPoints, redTeam.Points)
	}
	if blueTeam.Points != expBluePoints {
		t.Errorf("TestUpdateStatus(4) failed, expected Blue Team Points %d, got Blue Team Points %d", expBluePoints, blueTeam.Points)
	}
}
