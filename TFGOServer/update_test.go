package main

// update_test.go: tests for update functions in update.go

import (
	"testing"
)

func TestDistance(t *testing.T) {
	isTesting = true
	// same point
	dist := distance(Location{-5, 5}, Location{-5, 5})
	if !isAcceptableError(dist, 0, EPSILON) {
		t.Errorf("TestDistance(1) failed, expected distance 0, got distance %d", dist)
	}

	// different points
	dist = distance(Location{-4, -5}, Location{6, 5})
	if !isAcceptableError(dist, 14.1421356, EPSILON) {
		t.Errorf("TestDistance(2) failed, expected distance 14.1421356, got distance %d", dist)
	}
}

func TestInRange(t *testing.T) {
	isTesting = true
	// in range
	if inRange(Location{1, 0}, Location{0, 0}, 2) == false {
		t.Errorf("TestInRange(1) failed, expected TRUE, got FALSE")
	}

	// on border
	if inRange(Location{2, 0}, Location{0, 0}, 2) == false {
		t.Errorf("TestInRange(2) failed, expected TRUE, got FALSE")
	}

	// out of range
	if inRange(Location{3, 0}, Location{0, 0}, 2) == true {
		t.Errorf("TestInRange(3) failed, expected FALSE, got TRUE")
	}
}

func TestInGameBounds(t *testing.T) {
	isTesting = true
	g := makeSampleGame()

	// in bounds
	if inGameBounds(g, Location{50, 50}) == false {
		t.Errorf("TestInGameBounds(1) failed, expected TRUE, got FALSE")
	}

	// on corner
	if inGameBounds(g, Location{0, 0}) == false {
		t.Errorf("TestInGameBounds(2) failed, expected TRUE, got FALSE")
	}

	// on edge
	if inGameBounds(g, Location{0,50}) == false {
		t.Errorf("TestInGameBounds(3) failed, expected TRUE, got FALSE")
	}

	// out of bounds
	if inGameBounds(g, Location{420, 420}) == true {
		t.Errorf("TestInGameBounds(4) failed, expected FALSE, got TRUE")
	}
}

func TestUpdateLocation(t *testing.T) {
	isTesting = true
	g := makeSampleGame()
	cp := g.ControlPoints["CP2"]
	oliver := getOliver(g)

	// player moves out of bounds
	oliver.updateLocation(g, Location{420, 420}.locationToDegrees(), 0)
	if oliver.Status != OUTOFBOUNDS {
		t.Errorf("TestUpdateLocation(1) failed, expected Status OUTOFBOUNDS, got Status %s", playerStatusToString[oliver.Status])
	}

	// player enters control point
	expBlueCount := cp.BlueCount + 1
	oliver.updateLocation(g, cp.Location.locationToDegrees(), 0)
	if cp.BlueCount != expBlueCount {
		t.Errorf("TestUpdateLocation(2) failed, expected BlueCount %d, got BlueCount %d", expBlueCount, cp.BlueCount)
	}

	// player exits control point
	expBlueCount = cp.BlueCount - 1
	oliver.updateLocation(g, Location{0, 0}.locationToDegrees(), 0)
	if cp.BlueCount != expBlueCount {
		t.Errorf("TestUpdateLocation(3) failed, expected BlueCount %d, got BlueCount %d", expBlueCount, cp.BlueCount)
	}
}

func TestMovePayload(t *testing.T) {
	isTesting = true
	g := makeSampleGame()
	g.ControlPoints = nil
	g.ControlPoints = make(map[string]*ControlPoint)
	g.Mode = PAYLOAD
	g.generateObjectives(1)

	g.ControlPoints["Payload"].BlueCount = 0
	g.ControlPoints["Payload"].RedCount = 0
	movePayload(g)
	cpLoc := g.ControlPoints["Payload"].Location
	if cpLoc.X != 50.0 && cpLoc.Y != 50.0 {
		t.Errorf("TestMovePayload(1) failed, expected payload location to be (50,50), (%f,%f)", cpLoc.X, cpLoc.Y)
	}

	g.ControlPoints["Payload"].BlueCount = 1
	movePayload(g)
	cpLoc = g.ControlPoints["Payload"].Location
	if cpLoc.X != 50.0 && cpLoc.Y != 49.0 {
		t.Errorf("TestMovePayload(2) failed, expected payload location to be (50,49), got (%f,%f)", cpLoc.X, cpLoc.Y)
	}

	g.ControlPoints["Payload"].BlueCount = 0
	g.ControlPoints["Payload"].RedCount = 1
	movePayload(g)
	cpLoc = g.ControlPoints["Payload"].Location
	if cpLoc.X != 50.0 && cpLoc.Y != 50.0 {
		t.Errorf("TestMovePayload(3) failed, expected payload location to be (50,50), got (%f,%f)", cpLoc.X, cpLoc.Y)
	}
}

func TestUpdateStatus(t *testing.T) {
	isTesting = true
	g := makeSampleGame()
	cp1 := g.ControlPoints["CP1"]
	cp2 := g.ControlPoints["CP2"]
	redTeam := g.RedTeam
	blueTeam := g.BlueTeam
	oliver := getOliver(g)
	jenny := getJenny(g)
	anders := getAnders(g)
	brad := getBrad(g)

	// blue players exceed red by one; check capture progress
	expCaptureProg := cp2.CaptureProgress + 1
	//oliver.updateLocation(g, cp2.Location.locationToDegrees(), 0)
	cp2.updateStatus(g)
	if cp2.CaptureProgress != expCaptureProg {
		t.Errorf("TestUpdateStatus(1) failed, expected CaptureProgress %d, got CaptureProgress %d", expCaptureProg, cp2.CaptureProgress)
	}

	// equal players from each team; check capture progress
	oliver.updateLocation(g, Location{0, 0}.locationToDegrees(), 0)
	expCaptureProg = cp2.CaptureProgress
	cp2.updateStatus(g)
	if cp2.CaptureProgress != expCaptureProg {
		t.Errorf("TestUpdateStatus(2) failed, expected CaptureProgress %d, got CaptureProgress %d", expCaptureProg, cp2.CaptureProgress)
	}

	// both control points uncontrolled; check team points
	oliver.updateLocation(g, Location{0, 0}.locationToDegrees(), 0)
	jenny.updateLocation(g, Location{0, 0}.locationToDegrees(), 0)
	brad.updateLocation(g, Location{0, 0}.locationToDegrees(), 0)
	anders.updateLocation(g, Location{0, 0}.locationToDegrees(), 0)
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
	cp1.ControllingTeam = blueTeam
	cp2.ControllingTeam = blueTeam
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
