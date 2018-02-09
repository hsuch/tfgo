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
    10,
}

var loc3 = Location{
    10,
    0,
}

var loc4 = Location{
    10,
    10,
}

var loc5 = Location{
    5,
    5,
}

var loc6 = Location{
    42,
    17,
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
var inventory2 = [0]Pickup{}

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

}

func TestTakeHit(t *testing.T) {

}

func TestFire(t *testing.T) {

}

func TestHandleLoc(t *testing.T) {

}

func TestCanHit(t *testing.T) {

}
