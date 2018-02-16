package main

// setup_test.go: tests for game setup functions in setup.go

import (
	"testing"
	"time"
)

func TestFindCenter(t *testing.T) {
	isTesting = true
	g := makeSampleGame()
	x := g.findCenter().X
	y := g.findCenter().Y
	if !isAcceptableError(x, 50, EPSILON) && isAcceptableError(y, 50, EPSILON) {
		t.Errorf("TestFindCenter(1) failed, expected Location (50, 50), got Location (%d, %d)", x, y)
	}
}

func TestCreatePlayer(t *testing.T) {
	isTesting = true
	p := createPlayer(nil, "Alice", "testIcon")
	checkPlayerVitals(t, p, 100, 0, NORMAL, "TestCreatePlayer", "Alice")
	if p.Icon != "testIcon" {
		t.Errorf("TestCreatePlayer(1) failed, expected Icon testIcon, got Icon %s", p.Icon)
	}
}

func TestSetBoundaries(t *testing.T) {
	isTesting = true
	boundaries := []interface{} {
		map[string]interface{} {"X": 0.0, "Y": 0.0},
		map[string]interface{} {"X": meterToDegree(100.0), "Y": 0.0},
		map[string]interface{} {"X": meterToDegree(100.0), "Y": meterToDegree(100.0)},
		map[string]interface{} {"X": 0.0, "Y": meterToDegree(100.0)},
	}
	borders := []Border {
		Border{Location{0.0,0.0}, Direction{100.0,0.0}},
		Border{Location{100.0,0.0}, Direction{0.0,100.0}},
		Border{Location{100.0,100.0}, Direction{-100.0,0.0}},
		Border{Location{0.0,100.0}, Direction{0.0,-100.0}},
	}
	g := &Game{Name: "TestGame"}
	g.setBoundaries(boundaries)
	if len(g.Boundaries) != 4 {
		t.Errorf("TestSetBoundaries(1) failed, expected 4 borders, got %d", len(g.Boundaries))
	}
	for i, v := range borders {
		if g.Boundaries[i] != v {
			t.Errorf("TestSetBoundaries(2.%d) failed, expected Border{{%f, %f} {%f, %f}}, got Border{{%f, %f} {%f, %f}}",
				i, v.P.X, v.P.Y, v.D.X, v.D.Y,
				g.Boundaries[i].P.X, g.Boundaries[i].P.Y, g.Boundaries[i].D.X, g.Boundaries[i].D.Y)
		}
	}
}

func TestCreateGame(t *testing.T) {
	isTesting = true
	borders := []Border {
		Border{Location{0.0,0.0}, Direction{100.0,0.0}},
		Border{Location{100.0,0.0}, Direction{0.0,100.0}},
		Border{Location{100.0,100.0}, Direction{-100.0,0.0}},
		Border{Location{0.0,100.0}, Direction{0.0,-100.0}},
	}
	data := map[string]interface{} {
		"Name" : "G1Name",
		"Password" : "G1Password",
		"Description" : "G1Description",
		"PlayerLimit" : 8.0,
		"PointLimit" : 42.0,
		"TimeLimit" : "60m",
		"Mode" : "SingleCapture",
		"Boundaries" : []interface{} {
			map[string]interface{} {"X": 0.0, "Y": 0.0},
			map[string]interface{} {"X": meterToDegree(100.0), "Y": 0.0},
			map[string]interface{} {"X": meterToDegree(100.0), "Y": meterToDegree(100.0)},
			map[string]interface{} {"X": 0.0, "Y": meterToDegree(100.0)},
		},
		"Host" : map[string]interface{} {
			"Name" : "P1Name",
			"Icon" : "P1Icon",
		},

	}
	g, p := createGame(nil, data)

	if g.Name != data["Name"] {
		t.Errorf("TestCreateGame(1) failed, expected (Game) Name %s, got (Game) Name %s", data["Name"], g.Name)
	}
	if g.Password != data["Password"] {
		t.Errorf("TestCreateGame(2) failed, expected Password %s, got Password %s", data["Password"], g.Password)
	}
	if g.Description != data["Description"] {
		t.Errorf("TestCreateGame(3) failed, expected Description %s, got Description %s", data["Description"], g.Description)
	}
	if g.PlayerLimit != 8 {
		t.Errorf("TestCreateGame(4) failed, expected PlayerLimit %d, got PlayerLimit %d", 8, g.PlayerLimit)
	}
	if g.PointLimit != 42 {
		t.Errorf("TestCreateGame(5) failed, expected PointLimit %d, got PointLimit %d", 42, g.PointLimit)
	}
	if g.TimeLimit != 60 * time.Minute {
		t.Errorf("TestCreateGame(6) failed, expected TimeLimit %v, got TimeLimit %v", 60 * time.Minute, g.TimeLimit)
	}
	if g.Status != CREATING {
		t.Errorf("TestCreateGame(7) failed, expected Status %s, got Status %d", "CREATING", g.Status)
	}
	if g.Mode != SINGLECAP {
		t.Errorf("TestCreateGame(8) failed, expected Mode %s, got Mode %s", "SINGLECAP", g.Mode)
	}
	for i, v := range borders {
		if g.Boundaries[i] != v {
			t.Errorf("TestCreateGame(9.%d) failed, expected Border{{%d, %d} {%d, %d}}, got Border{{%d, %d} {%d, %d}}",
				i, v.P.X, v.P.Y, v.D.X, v.D.Y,
				g.Boundaries[i].P.X, g.Boundaries[i].P.Y, g.Boundaries[i].D.X, g.Boundaries[i].D.Y)
		}
	}
	if p.Name != data["Host"].(map[string]interface{})["Name"] {
		t.Errorf("TestCreateGame(10) failed, expected (Player) Name %s, got (Player) Name %s",
			data["Host"].(map[string]interface{})["Name"], p.Name)
	}
	if p.Icon != data["Host"].(map[string]interface{})["Icon"] {
		t.Errorf("TestCreateGame(11) failed, expected Icon %s, got Icon %s",
			data["Host"].(map[string]interface{})["Icon"], p.Icon)
	}
	if g.RedTeam.Name != "Red" {
		t.Errorf("TestCreateGame(12) failed, expected (RedTeam) Name %s, got (RedTeam) Name %s", "Red", g.RedTeam.Name)
	}
	if g.BlueTeam.Name != "Blue" {
		t.Errorf("TestCreateGame(13) failed, expected (BlueTeam) Name %s, got (BlueTeam) Name %s", "Blue", g.BlueTeam.Name)
	}
	if g.Players[data["Host"].(map[string]interface{})["Name"].(string)] != p {
		t.Errorf("TestCreateGame(14) failed, expected Player %s to be in Game Player List, got Player %s",
			p.Name, g.Players[data["Host"].(map[string]interface{})["Name"].(string)])
	}
}

func TestGenerateObjectives(t *testing.T) {
	isTesting = true
	// a game arena that is a wide rectangle
	// 1 control point
	g1 := &Game{
		Boundaries: []Border{
			{Location{0, 0}, Direction{200, 0}},
			{Location{200, 0}, Direction{0, 100}},
			{Location{200, 100}, Direction{-200, 0}},
			{Location{0, 100}, Direction{0, -100}},
		},
		RedTeam: &Team{Name:"RedTeam"},
		BlueTeam: &Team{Name:"BlueTeam"},
	}
	g1.generateObjectives(1)
	if (g1.RedTeam.Base != Location{195.0, 50.0}) {
		t.Errorf("TestGenerateObjectives(1) failed, expected Location{195,50}, got Location{%f,%f}", g1.RedTeam.Base.X, g1.RedTeam.Base.Y)
	}
	if g1.RedTeam.BaseRadius != 3.0 {
		t.Errorf("TestGenerateObjectives(2) failed, expected BaseRadius 3, got BaseRadius %f", g1.RedTeam.BaseRadius)
	}
	if (g1.BlueTeam.Base != Location{5.0, 50.0}) {
		t.Errorf("TestGenerateObjectives(3) failed, expected Location{5,50}, got Location{%f,%f}", g1.BlueTeam.Base.X, g1.BlueTeam.Base.Y)
	}
	if g1.BlueTeam.BaseRadius != 3.0 {
		t.Errorf("TestGenerateObjectives(4) failed, expected BaseRadius 3, got BaseRadius %f", g1.BlueTeam.BaseRadius)
	}
	if len(g1.ControlPoints) != 1 {
		t.Errorf("TestGenerateObjectives(5) failed, expected 1 ControlPoint, got %d", len(g1.ControlPoints))
	}
	if !inGameBounds(g1, g1.ControlPoints["CP1"].Location) {
		t.Errorf("TestGenerateObjectives(6) failed, expected inGameBounds(CP) to be TRUE, got FALSE")
	}

	// a game arena that is a tall rectangle
	// 2 control points
	g2 := &Game{
		Boundaries: []Border{
			{Location{0, 0}, Direction{100, 0}},
			{Location{100, 0}, Direction{0, 200}},
			{Location{100, 200}, Direction{-100, 0}},
			{Location{0, 200}, Direction{0, -200}},
		},
		RedTeam: &Team{Name:"RedTeam"},
		BlueTeam: &Team{Name:"BlueTeam"},
	}
	g2.generateObjectives(2)
	if (g2.RedTeam.Base != Location{50.0, 195.0}) {
		t.Errorf("TestGenerateObjectives(7) failed, expected Location{50, 195}, got Location{%f,%f}", g2.RedTeam.Base.X, g2.RedTeam.Base.Y)
	}
	if g2.RedTeam.BaseRadius != 3.0 {
		t.Errorf("TestGenerateObjectives(8) failed, expected BaseRadius 3, got BaseRadius %f", g2.RedTeam.BaseRadius)
	}
	if (g2.BlueTeam.Base != Location{50.0, 5.0}) {
		t.Errorf("TestGenerateObjectives(9) failed, expected Location{50,5}, got Location{%f,%f}", g2.BlueTeam.Base.X, g2.BlueTeam.Base.Y)
	}
	if g2.BlueTeam.BaseRadius != 3.0 {
		t.Errorf("TestGenerateObjectives(10) failed, expected BaseRadius 3, got BaseRadius %f", g2.BlueTeam.BaseRadius)
	}
	if len(g2.ControlPoints) != 2 {
		t.Errorf("TestGenerateObjectives(11) failed, expected 1 ControlPoint, got %d", len(g2.ControlPoints))
	}
	if !inGameBounds(g2, g2.ControlPoints["CP1"].Location) {
		t.Errorf("TestGenerateObjectives(12) failed, expected inGameBounds(CP) to be TRUE, got FALSE")
	}
	if !inGameBounds(g2, g2.ControlPoints["CP2"].Location) {
		t.Errorf("TestGenerateObjectives(13) failed, expected inGameBounds(CP) to be TRUE, got FALSE")
	}
}
