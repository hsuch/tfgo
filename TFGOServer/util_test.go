package main

var boundaries = [4]Location {
	{0, 0},
	{100, 0},
	{100, 100},
	{0, 100},
}

var jenny = Player {
	Name: "Jenny",
	Team: &redTeam,
	Status: NORMAL,
	Health: 100,
	Armor: 0,
	Weapon: SWORD,
	Location: boundaries[2],
}

var brad = Player {
	Name: "Brad",
	Team: &redTeam,
	Status: NORMAL,
	Health: 80,
	Armor: 30,
	Weapon: SWORD,
	Location: boundaries[3],
}

var anders = Player {
	Name: "Anders",
	Team: &blueTeam,
	Status: NORMAL,
	Health: 10,
	Armor: 5,
	Weapon: SWORD,
	Location: boundaries[2],
}

var oliver = Player {
	Name: "Oliver",
	Team: &blueTeam,
	Status: NORMAL,
	Health: 104,
	Armor: 2,
	Weapon: SWORD,
	Location: boundaries[1],
}

var redTeam = Team {
	Name: "Red Team",
	Base: boundaries[0],
	Points: 5,
	Players: []*Player{&jenny, &brad},
}

var blueTeam = Team {
	Name: "Blue Team",
	Base: boundaries[1],
	Points: 6,
	Players: []*Player{&anders, &oliver},
}

var cp1 = ControlPoint {
	ID: "CP1",
	Location: boundaries[3],
	Radius: 5,
	RedCount: 1,
	BlueCount: 0,
}

var cp2 = ControlPoint {
	ID: "CP2",
	Location: boundaries[2],
	Radius: 5,
	RedCount: 1,
	BlueCount: 1,
}

var testGame = Game {
	ID: "G1",
	Name: "Test Game",
	Password: "tfgo",
	Status: PLAYING,
	Mode: SINGLECAP,
	RedTeam: &redTeam,
	BlueTeam: &blueTeam,
	Boundaries: boundaries,
	ControlPoints: map[string]*ControlPoint {
		"CP1" : &cp1,
		"CP2" : &cp2,
	},
}
