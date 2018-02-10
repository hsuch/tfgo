package TFGOServer

import (
    "testing"
    "time"
    "math"
)

var loc1 = Location{
    0,
    0,
}

var loc2 = Location{
    0,
    100,
}

var loc3 = Location{
    100,
    0,
}

var loc4 = Location{
    100,
    100,
}

var loc5 = Location{
    50,
    50,
}

var loc6 = Location{
    420,
    170,
}

var direction1 = Direction{
    0,
    1,
}

var direction2 = Direction{
    -1,
    8,
}

var direction3 = Direction{
    1,
    8,
}

var direction4 = Direction{
    -5,
    0,
}

var direction5 = Direction{
    5,
    0,
}

var direction6 = Direction{
    1,
    12,
}

var weapon1 = Weapon{
	Name: "TestWeapon",
	Damage: 25,
	Spread: math.Pi/2,
	Range: 10,
	ClipSize: 42,
	ShotReload: time.Second,
	ClipReload: time.Second * 5,
}

var jenny = Player{
    Name: "Jenny",
    Team: redTeam,
    Status: 0,
    Health: 100,
    Armor: 0,
    Weapon: SWORD,
    Inventory: inventory1,
    Location: loc4,
}

var brad = Player{
    Name: "Brad",
    Team: redTeam,
    Status: 0,
    Health: 80,
    Armor: 20,
    Weapon: SWORD,
    Inventory: inventory2,
    Location: loc2,
}

var anders = Player{
    Name: "Anders",
    Team: blueTeam,
    Status: 0,
    Health: 90,
    Armor: 5,
    Weapon: SWORD,
    Inventory: inventory1,
    Location: loc4,
}

var oliver = Player{
    Name: "Oliver",
    Team: blueTeam,
    Status: 1,
    StatusTimer: time.NewTimer(30 * time.second),
    Health: 104,
    Armor: 2,
    Weapon: SWORD,
    Inventory: inventory2,
    Location: loc5,
}

var redTeam = Team{
    Name: "Red Team",
    Base: loc1,
    Points: 5,
    Players: [2]Player{jenny, brad},
}

var blueTeam = Team{
    Name: "Blue Team",
    Base: loc3,
    Points: 6,
    Players: [2]Player{anders, oliver},
}

var emptyTeam = Team{
    Name: "Lonely Team",
    Base: loc2,
    Points: 0,
    Players: nil,
}

var cp1 = ControlPoint{
    ID: "CP1",
    Location: loc2,
    Radius: 2,
    PayloadPath: [2]Location{loc1, loc2},
    PayloadLoc: loc1,
    RedCount: 0,
    BlueCount: 0,
    ControllingTeam: nil,
    CaptureProgress: 0,
}

var cp2 = ControlPoint{
    ID: "CP2",
    Location: loc4,
    Radius: 2.1,
    PayloadPath: [2]Location{loc3, loc4},
    PayloadLoc: loc3,
    RedCount: 1,
    BlueCount: 0,
    ControllingTeam: nil,
    CaptureProgress: 0,
}

var inventory1 = [1]Pickup{SWORD}
var inventory2 = nil

var testGame = Game{
    ID: "G1",
    Name: "Test Game",
    Password: "tfgo",
    Status: 0,
    Mode: 0,
    Timer: time.NewTimer(time.hour),
    RedTeam: redTeam,
    BlueTeam: blueTeam,
    Boundaries: [4]Location{loc1, loc2, loc3, loc4},
    ControlPoints: map[string]ControlPoint{
        "CP1" : cp1,
        "CP2" : cp2,
    },
}

func TestGetPlayerLocs(t *testing.T) {
    // testing with a team with 2 members
    teamLocs := GetTeamLocs(redTeam)
    if teamLocs[0] != loc4 {
        t.Errorf("TestGetTeamLocs(redTeam) failed, got: (%f,%f), want: (100,100).", teamLocs[0].X, teamLocs[0].Y)
    }
    if teamLocs[1] != loc2 {
        t.Errorf("TestGetTeamLocs(redTeam) failed, got: (%f,%f), want: (0,100).", teamLocs[1].X, teamLocs[1].Y)
    }

    // testing with a team with no members
    teamLocs = GetTeamLocs(emptyTeam)
    if teamLocs != nil {
        t.Errorf("TestGetTeamLocs(emptyTeam) failed, expected output length 0, got length %d.", len(teamLocs))
    }
}

func TestTakeHit(t *testing.T) {

}

func TestFire(t *testing.T) {

}

func TestHandleLoc(t *testing.T) {
    // 1. Testing a player out of bounds
    (&oliver).handleLoc(loc6)
    if (oliver.Status != OUTOFBOUNDS) {
        t.Errorf("TestHandleLoc(1) failed, expected Status OUTOFBOUNDS (%d), got Status %d", OUTOFBOUNDS, oliver.Status)
        // test if timer is running?
    }

    // 2. Testing for player count addition to a control point
    expBlueCount := cp2.BlueCount + 1
    (&oliver).handleLoc(loc4)
    if (cp2.BlueCount != expBlueCount) {
        t.Errorf("TestHandleLoc(2) failed, expected BlueCount %d, got BlueCount %d", expBlueCount, cp2.BlueCount)
    }

    // 3. Testing for player count subtraction from a control point
    expBlueCount = cp2.BlueCount - 1
    (&oliver).handleLoc(loc5)
    if (cp2.BlueCount != expBlueCount) {
        t.Errorf("TestHandleLoc(3) failed, expected BlueCount %d, got BlueCount %d", expBlueCount, cp2.BlueCount)
    }
}

func TestUpdateStatus(t *testing.T) {
    /* Assumptions:
     * - handleLoc executes properly
     * - CaptureProgress is + for Blue Team, - for Red Team
     */

    // 1. Testing control point contest when 2:1 in favor of Blue Team (Oliver and Anders (blue) vs. Jenny (red) at loc4)
    expCaptureProg := cp2.CaptureProgress + 1
    (&oliver).handleLoc(loc4)
    cp2.updateStatus()
    if (cp2.CaptureProgress != expCaptureProg) {
        t.Errorf("TestUpdateStatus(1) failed, expected CaptureProgress %d, got CaptureProgress %d", expCaptureProg, cp2.CaptureProgress)
    }

    // 2. Testing control point contest when 1:1 (Anders (blue) vs. Jenny (red) at loc4)
    expCaptureProg = cp2.CaptureProgress
    (&oliver).handleLoc(loc5)
    cp2.updateStatus()
    if (cp2.CaptureProgress != expCaptureProg) {
        t.Errorf("TestUpdateStatus(2) failed, expected CaptureProgress %d, got CaptureProgress %d", expCaptureProg, cp2.CaptureProgress)
    }

    // 3. Testing team point increase with no capture points
    cp1.ControllingTeam = nil
    cp2.ControllingTeam = nil
    expRedPoints := redTeam.Points
    expBluePoints := blueTeam.Points
    cp1.updateStatus()
    cp2.updateStatus()
    if (redTeam.Points != expRedPoints) {
        t.Errorf("TestUpdateStatus(3) failed, expected Red Team Points %d, got Red Team Points %d", expRedPoints, redTeam.Points)
    }
    if (blueTeam.Points != expBluePoints) {
        t.Errorf("TestUpdateStatus(3) failed, expected Blue Team Points %d, got Blue Team Points %d", expBluePoints, blueTeam.Points)
    }

    // 4. Testing team point increase with 2 capture points
    cp1.ControllingTeam = blueTeam
    cp2.ControllingTeam = blueTeam
    expRedPoints = redTeam.Points
    expBluePoints = blueTeam.Points + 2
    cp1.updateStatus()
    cp2.updateStatus()
    if (redTeam.Points != expRedPoints) {
        t.Errorf("TestUpdateStatus(4) failed, expected Red Team Points %d, got Red Team Points %d", expRedPoints, redTeam.Points)
    }
    if (blueTeam.Points != expBluePoints) {
        t.Errorf("TestUpdateStatus(4) failed, expected Blue Team Points %d, got Blue Team Points %d", expBluePoints, blueTeam.Points)
    }
}

func TestCanHit(t *testing.T) {
    // if the shot is to the left of the target but within the spread
    success := CanHit(direction1, direction2, weapon1)
    if !success {
        t.Errorf("TestCanHit(1) failed, expected TRUE but got FALSE.")
    }

    // if the shot is to the right of the target but within the spread
    success = CanHit(direction1, direction3, weapon1)
    if !success {
        t.Errorf("TestCanHit(2) failed, expected TRUE but got FALSE.")
    }

    // if the shot is to the left of the target and outside the spread
    success = CanHit(direction1, direction4, weapon1)
    if success {
        t.Errorf("TestCanHit(3) failed, expected FALSE but got TRUE.")
    }

    // if the shot is to the right of the target and outside the spread
    success = CanHit(direction1, direction5, weapon1)
    if success {
        t.Errorf("TestCanHit(4) failed, expected FALSE but got TRUE.")
    }

    // if the target is not within range
    success = CanHit(direction1, direction6, weapon1)
    if success {
        t.Errorf("TestCanHit(5) failed, expected FALSE but got TRUE.")
    }
}
