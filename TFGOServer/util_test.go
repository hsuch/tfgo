package main

// utility functions used for testing, including:
// (1) checking if a value is within a specified error
// (2) creating instances of a sample game and its components

import "math"

const EPSILON = 0.05

func isAcceptableError(testValue float64, expectedValue float64, errorThreshold float64) bool {
	if expectedValue == 0 {
		return math.Abs(testValue) < 0.0001 // four decimal points should be sufficient for our purposes
	}

	error := math.Abs(expectedValue - testValue) / expectedValue
	return error <= errorThreshold
}

func makeJenny(team *Team) *Player {
	return &Player {
		Name: "Jenny",
		Team: team,
		Status: NORMAL,
		Health: 100,
		Armor: 0,
		Location: Location{49, 75},
	}
}

func makeBrad(team *Team) *Player {
	return &Player {
		Name: "Brad",
		Team: team,
		Status: NORMAL,
		Health: 80,
		Armor: 30,
		Location: Location{49, 24},
	}
}

func makeAnders(team *Team) *Player {
	return &Player {
		Name: "Anders",
		Team: team,
		Status: NORMAL,
		Health: 10,
		Armor: 5,
		Location: Location{49.5, 75},
	}
}

func makeOliver(team *Team) *Player {
	return &Player {
		Name: "Oliver",
		Team: team,
		Status: NORMAL,
		Health: 95,
		Armor: 10,
		Location: Location{50, 75},
	}
}

func makeRedTeam() *Team {
	return &Team {
		Name: "Red Team",
		Base: Location{25, 50},
		Points: 98,
	}
}

func makeBlueTeam() *Team {
	return &Team {
		Name: "Blue Team",
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
		ControllingTeam: nil,
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
		ControllingTeam: nil,
		CaptureProgress: -3,
	}
}

func makeSampleGame() *Game {
	redTeam := makeRedTeam()
	blueTeam := makeBlueTeam()
	return &Game {
		ID: "G1",
		Name: "Test Game",
		Password: "tfgo",
		Status: PLAYING,
		Mode: MULTICAP,
		RedTeam: redTeam,
		BlueTeam: blueTeam,
		Players: map[string]*Player{
			"jenny" : makeJenny(redTeam),
			"brad" : makeBrad(redTeam),
			"oliver" : makeOliver(blueTeam),
			"anders" : makeAnders(blueTeam),
		},
		Boundaries: []Border{
			{Location{0, 0}, Direction{1, 0}},
			{Location{100, 0}, Direction{0, 1}},
			{Location{100, 100}, Direction{-1, 0}},
			{Location{0, 100}, Direction{0, -1}},
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
	return game.Players["jenny"]
}

func getOliver(game *Game) *Player {
	return game.Players["oliver"]
}

func getAnders(game *Game) *Player {
	return game.Players["anders"]
}

func getBrad(game *Game) *Player {
	return game.Players["brad"]
}