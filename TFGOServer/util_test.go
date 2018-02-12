package main

// functions for creating instances of a sample game and its components,
// used for testing.

func makeJenny() *Player {
	return &Player {
		Name: "Jenny",
		Team: RED,
		Status: NORMAL,
		Health: 100,
		Armor: 0,
		Location: Location{49, 75},
	}
}

func makeBrad() *Player {
	return &Player {
		Name: "Brad",
		Team: RED,
		Status: NORMAL,
		Health: 80,
		Armor: 30,
		Location: Location{49, 24},
	}
}

func makeAnders() *Player {
	return &Player {
		Name: "Anders",
		Team: BLUE,
		Status: NORMAL,
		Health: 10,
		Armor: 5,
		Location: Location{49.5, 75},
	}
}

func makeOliver() *Player {
	return &Player {
		Name: "Oliver",
		Team: BLUE,
		Status: NORMAL,
		Health: 95,
		Armor: 10,
		Location: Location{50, 75},
	}
}

func makeRedTeam() *Team {
	return &Team {
		Name: "Red Team",
		Players: map[string]*Player{
			"jenny" : makeJenny(),
			"brad" : makeBrad(),
		},
		Base: Location{25, 50},
		Points: 98,
	}
}

func makeBlueTeam() *Team {
	return &Team {
		Name: "Blue Team",
		Players: map[string]*Player{
			"oliver" : makeOliver(),
			"anders" : makeAnders(),
		},
		Base: Location{75, 50},
		Points: 72,
	}
}

func makeCP1() *ControlPoint {
	return &ControlPoint {
		ID: "CP1",
		Location: Location{50, 25},
		Radius: 5,
		RedCount: 1,
		BlueCount: 0,
		ControllingTeam: NEUTRAL,
		CaptureProgress: -6,
	}
}

func makeCP2() *ControlPoint {
	return &ControlPoint {
		ID: "CP2",
		Location: Location{50, 75},
		Radius: 5,
		RedCount: 1,
		BlueCount: 2,
		ControllingTeam: RED,
		CaptureProgress: -3,
	}
}

func makeSampleGame() *Game {
	return &Game {
		ID: "G1",
		Name: "Test Game",
		Password: "tfgo",
		Status: PLAYING,
		Mode: MULTICAP,
		RedTeam: makeRedTeam(),
		BlueTeam: makeBlueTeam(),
		Boundaries: [4]Location{
			{0, 0},
			{100, 0},
			{0, 100},
			{100, 100},
		},
		ControlPoints: map[string]*ControlPoint {
			"CP1" : makeCP1(),
			"CP2" : makeCP2(),
		},
	}
}

// functions used to retrieve a specific player from the sample game
// defined above. used for testing.

func getJenny(game *Game) *Player {
	return game.RedTeam.Players["jenny"]
}

func getOliver(game *Game) *Player {
	return game.BlueTeam.Players["oliver"]
}

func getAnders(game *Game) *Player {
	return game.BlueTeam.Players["anders"]
}

func getBrad(game *Game) *Player {
	return game.RedTeam.Players["brad"]
}