package TFGOServer

import (
    "testing"
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
    Location: loc5,
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
    Location: loc6,
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
    Location: loc1,
    Radius: 2,
    PayloadPath: [2]Location{loc1, loc2},
    PayloadLoc: loc1,
    RedCount: 0,
    BlueCount: 0,
    CaptureStatus: 0,
}

var cp2 = ControlPoint{
    ID: "CP2",
    Location: loc4,
    Radius: 2.1,
    PayloadPath: [2]Location{loc3, loc4},
    PayloadLoc: loc3,
    RedCount: 1,
    BlueCount: 0,
    CaptureStatus: 1,
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

func TestUpdateStatus(t *testing.T) {
    /* The following is an example of a unit test: */
    // out := Function()
    // if out != expected_out {
    //    t.Errorf("Function test failed, got: %d, want: %d.", out, expected_out)
    // }
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
