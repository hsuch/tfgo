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
    Points: 5,
    Players: [2]Player{jenny, brad},
}

var blueTeam = Team{
    Name: "Blue Team",
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
    RedBase: loc1,
    BlueTeam: blueTeam,
    BlueBase: loc3,
    Boundaries: [4]Location{loc1, loc2, loc3, loc4},
    ControlPoints: map[string]ControlPoint{
        "CP1" : cp1,
        "CP2" : cp2,
    },
}

func TestLocationSetGet(*testing.T)
{

}

func TestUpdateStatus(*testing.T)
{

}

func TestGetTeamLoc(*testing.T)
{

}

func TestSetupGame(*testing.T)
{

}

func TestAddPlayer(*testing.T)
{

}

func TestSetWinner(*testing.T)
{

}

func TestTakeDamage(*testing.T)
{

}

func TestHeal(*testing.T)
{

}

func TestShoot(*testing.T)
{

}

func TestCheckInventory(*testing.T)
{

}

func TestAddRemoveItem(*testing.T)
{

}

func TestHandlePlayerLoc(*testing.T)
{

}

func TestCanHit(*testing.T)
{

}

func TestNearestHit(*testing.T)
{

}

func TestPickupHandler(*testing.T)
{

}
