package main

import "testing"

func TestHandleLoc(t *testing.T) {
	p1 := oliver

	// 1. Testing a player out of bounds
	p1.handleLoc(Location{420, 170})
	if p1.Status != OUTOFBOUNDS {
		t.Errorf("TestHandleLoc(1) failed, expected Status OUTOFBOUNDS (%d), got Status %d", OUTOFBOUNDS, p1.Status)
	}

	// 2. Testing for player count addition to a control point
	expBlueCount := cp2.BlueCount + 1
	p1.handleLoc(boundaries[2])
	if cp2.BlueCount != expBlueCount {
		t.Errorf("TestHandleLoc(2) failed, expected BlueCount %d, got BlueCount %d", expBlueCount, cp2.BlueCount)
	}

	// 3. Testing for player count subtraction from a control point
	expBlueCount = cp2.BlueCount - 1
	p1.handleLoc(boundaries[1])
	if cp2.BlueCount != expBlueCount {
		t.Errorf("TestHandleLoc(3) failed, expected BlueCount %d, got BlueCount %d", expBlueCount, cp2.BlueCount)
	}
}

func TestUpdateStatus(t *testing.T) {
	/* Assumptions:
	 * - handleLoc executes properly
	 * - CaptureProgress is + for Blue Team, - for Red Team
	 */
	p1 := oliver

	// 1. Testing control point contest when 2:1 in favor of Blue Team (Oliver and Anders (blue) vs. Jenny (red) at loc4)
	expCaptureProg := cp2.CaptureProgress + 1
	p1.handleLoc(boundaries[2])
	cp2.updateStatus()
	if cp2.CaptureProgress != expCaptureProg {
		t.Errorf("TestUpdateStatus(1) failed, expected CaptureProgress %d, got CaptureProgress %d", expCaptureProg, cp2.CaptureProgress)
	}

	// 2. Testing control point contest when 1:1 (Anders (blue) vs. Jenny (red) at loc4)
	expCaptureProg = cp2.CaptureProgress
	p1.handleLoc(boundaries[1])
	cp2.updateStatus()
	if cp2.CaptureProgress != expCaptureProg {
		t.Errorf("TestUpdateStatus(2) failed, expected CaptureProgress %d, got CaptureProgress %d", expCaptureProg, cp2.CaptureProgress)
	}

	// 3. Testing team point increase with no capture points
	cp1.ControllingTeam = nil
	cp2.ControllingTeam = nil
	expRedPoints := redTeam.Points
	expBluePoints := blueTeam.Points
	cp1.updateStatus()
	cp2.updateStatus()
	if redTeam.Points != expRedPoints {
		t.Errorf("TestUpdateStatus(3) failed, expected Red Team Points %d, got Red Team Points %d", expRedPoints, redTeam.Points)
	}
	if blueTeam.Points != expBluePoints {
		t.Errorf("TestUpdateStatus(3) failed, expected Blue Team Points %d, got Blue Team Points %d", expBluePoints, blueTeam.Points)
	}

	// 4. Testing team point increase with 2 capture points
	cp1.ControllingTeam = &blueTeam
	cp2.ControllingTeam = &blueTeam
	expRedPoints = redTeam.Points
	expBluePoints = blueTeam.Points + 2
	cp1.updateStatus()
	cp2.updateStatus()
	if redTeam.Points != expRedPoints {
		t.Errorf("TestUpdateStatus(4) failed, expected Red Team Points %d, got Red Team Points %d", expRedPoints, redTeam.Points)
	}
	if blueTeam.Points != expBluePoints {
		t.Errorf("TestUpdateStatus(4) failed, expected Blue Team Points %d, got Blue Team Points %d", expBluePoints, blueTeam.Points)
	}
}