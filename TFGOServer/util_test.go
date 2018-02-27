package main

// utility functions used for testing, including:
// (1) checking if a value is within a specified error
// (2) creating instances of a sample game and its components

import (
	"math"
)

const EPSILON = 0.05

func isAcceptableError(testValue float64, expectedValue float64, errorThreshold float64) bool {
	if expectedValue == 0 {
		return math.Abs(testValue) < 0.0001 // four decimal points should be sufficient for our purposes
	}

	err := math.Abs(expectedValue - testValue) / expectedValue
	return err <= errorThreshold
}

func makeJenny(team *Team, cp *ControlPoint) *Player {
	jenny := Player {
		ID: "jenny",
		Name: "Jenny",
		Team: team,
		Chan: make(chan map[string]interface{}),
		Status: NORMAL,
		Health: 100,
		Armor: 0,
		Location: Location{49, 75},
		OccupyingPoint: cp,
	}
	go jenny.sender()
	return &jenny
}

func makeBrad(team *Team, cp *ControlPoint) *Player {
	brad := Player {
		ID: "brad",
		Name: "Brad",
		Team: team,
		Chan: make(chan map[string]interface{}),
		Status: NORMAL,
		Health: 80,
		Armor: 30,
		Location: Location{49, 24},
		OccupyingPoint: cp,
	}
	go brad.sender()
	return &brad
}

func makeAnders(team *Team, cp *ControlPoint) *Player {
	anders := Player {
		ID: "anders",
		Name: "Anders",
		Team: team,
		Chan: make(chan map[string]interface{}),
		Status: NORMAL,
		Health: 10,
		Armor: 5,
		Location: Location{49.5, 75},
		OccupyingPoint: cp,
	}
	go anders.sender()
	return &anders
}

func makeOliver(team *Team, cp *ControlPoint) *Player {
	oliver := Player {
		ID: "oliver",
		Name: "Oliver",
		Team: team,
		Chan: make(chan map[string]interface{}),
		Status: NORMAL,
		Health: 95,
		Armor: 10,
		Location: Location{50, 75},
		OccupyingPoint: cp,
	}
	go oliver.sender()
	return &oliver
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
	cp1 := makeCP1()
	cp2 := makeCP2()
	return &Game {
		HostID: "jenny",
		Name: "Test Game",
		Password: "tfgo",
		Status: PLAYING,
		Mode: MULTICAP,
		RedTeam: redTeam,
		BlueTeam: blueTeam,
		Players: map[string]*Player{
			"jenny" : makeJenny(redTeam, cp2),
			"brad" : makeBrad(redTeam, cp1),
			"oliver" : makeOliver(blueTeam, cp2),
			"anders" : makeAnders(blueTeam, cp2),
		},
		Boundaries: []Border{
			{Location{0, 0}, Direction{100, 0}},
			{Location{100, 0}, Direction{0, 100}},
			{Location{100, 100}, Direction{-100, 0}},
			{Location{0, 100}, Direction{0, -100}},
		},
		ControlPoints: map[string]*ControlPoint {
			"CP1" : cp1,
			"CP2" : cp2,
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
