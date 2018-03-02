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
		{Location{0.0,0.0}, Direction{100.0,0.0}},
		{Location{100.0,0.0}, Direction{0.0,100.0}},
		{Location{100.0,100.0}, Direction{-100.0,0.0}},
		{Location{0.0,100.0}, Direction{0.0,-100.0}},
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

	p0 := map[string]interface{} {"X": meterToDegree(0.0), "Y": meterToDegree(0.0)}
	p1 := map[string]interface{} {"X": meterToDegree(-10.0), "Y": meterToDegree(50.0)}
	p2 := map[string]interface{} {"X": meterToDegree(-10.0), "Y": meterToDegree(100.0)}
	p3 := map[string]interface{} {"X": meterToDegree(0.0), "Y": meterToDegree(120.0)}
	p4 := map[string]interface{} {"X": meterToDegree(50.0), "Y": meterToDegree(100.0)}
	p5 := map[string]interface{} {"X": meterToDegree(100.0), "Y": meterToDegree(80.0)}
	p6 := map[string]interface{} {"X": meterToDegree(100.0), "Y": meterToDegree(50.0)}
	p7 := map[string]interface{} {"X": meterToDegree(50.0), "Y": meterToDegree(25.0)}
	boundaries2 := []interface{} {p0,p2,p4,p6,p1,p3,p5,p7}
	borders2 := []Border {
		{Location{0.0,0.0}, Direction{100.0,50.0}},
		{Location{100.0,50.0}, Direction{-50.0,-25.0}},
		{Location{50.0,25.0}, Direction{50.0,55.0}},
		{Location{100.0,80.0}, Direction{-50.0,20.0}},
		{Location{50.0,100.0}, Direction{-50.0,20.0}},
		{Location{0.0,120.0}, Direction{-10.0,-20.0}},
		{Location{-10.0,100.0}, Direction{0.0,-50.0}},
		{Location{-10.0,50.0}, Direction{10.0,-50.0}},
	}

	g2 := &Game{Name: "TestGame"}
	g2.setBoundaries(boundaries2)
	if len(g2.Boundaries) != 8 {
		t.Errorf("TestSetBoundaries(3) failed, expected 8 borders, got %d", len(g2.Boundaries))
	}
	for i, v := range borders2 {
		if g2.Boundaries[i] != v {
			t.Errorf("TestSetBoundaries(4.%d) failed, expected Border{{%f, %f} {%f, %f}}, got Border{{%f, %f} {%f, %f}}",
				i, v.P.X, v.P.Y, v.D.X, v.D.Y,
				g2.Boundaries[i].P.X, g2.Boundaries[i].P.Y, g2.Boundaries[i].D.X, g2.Boundaries[i].D.Y)
		}
	}
}

func TestCreateGame(t *testing.T) {
	isTesting = true
	borders := []Border {
		{Location{0.0,0.0}, Direction{100.0,0.0}},
		{Location{100.0,0.0}, Direction{0.0,100.0}},
		{Location{100.0,100.0}, Direction{-100.0,0.0}},
		{Location{0.0,100.0}, Direction{0.0,-100.0}},
	}
	data := map[string]interface{} {
		"Name" : "G1Name",
		"Password" : "G1Password",
		"Description" : "G1Description",
		"PlayerLimit" : 8.0,
		"PointLimit" : 42.0,
		"TimeLimit" : 60.0,
		"Mode" : "SingleCapture",
		"Boundaries" : []interface{} {
			map[string]interface{} {"X": 0.0, "Y": 0.0},
			map[string]interface{} {"X": meterToDegree(100.0), "Y": 0.0},
			map[string]interface{} {"X": meterToDegree(100.0), "Y": meterToDegree(100.0)},
			map[string]interface{} {"X": 0.0, "Y": meterToDegree(100.0)},
		},
	}

	p := &Player{ID: "P1ID", Name: "P1Name", Icon: "P1Icon"}
	g := p.createGame(nil, data)

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
	if g.RedTeam.Name != "Red" {
		t.Errorf("TestCreateGame(12) failed, expected (RedTeam) Name %s, got (RedTeam) Name %s", "Red", g.RedTeam.Name)
	}
	if g.BlueTeam.Name != "Blue" {
		t.Errorf("TestCreateGame(13) failed, expected (BlueTeam) Name %s, got (BlueTeam) Name %s", "Blue", g.BlueTeam.Name)
	}
	if g.Players[p.ID] != p {
		t.Errorf("TestCreateGame(14) failed, expected Player %s to be in Game Player List, got Player %s",
			p.Name, g.Players[p.ID])
	}
}

func TestGenerateObjectives(t *testing.T) {
	isTesting = true
	// a game arena that is a wide rectangle
	// SINGLECAP mode (same setup as PAYLOAD mode)
	g1 := &Game{
		Boundaries: []Border{
			{Location{300, 300}, Direction{200, 0}},
			{Location{500, 300}, Direction{0, 100}},
			{Location{500, 400}, Direction{-200, 0}},
			{Location{300, 400}, Direction{0, -100}},
		},
		RedTeam: &Team{Name:"RedTeam"},
		BlueTeam: &Team{Name:"BlueTeam"},
		Mode: SINGLECAP,
	}
	g1.generateObjectives(1)
	time.Sleep(200 * time.Millisecond)
	if (g1.RedTeam.Base != Location{491.0, 350.0}) {
		t.Errorf("TestGenerateObjectives(1) failed, expected Location{491,350}, got Location{%f,%f}", g1.RedTeam.Base.X, g1.RedTeam.Base.Y)
	}
	if g1.RedTeam.BaseRadius != 7.0 {
		t.Errorf("TestGenerateObjectives(2) failed, expected BaseRadius 7, got BaseRadius %f", g1.RedTeam.BaseRadius)
	}
	if (g1.BlueTeam.Base != Location{309.0, 350.0}) {
		t.Errorf("TestGenerateObjectives(3) failed, expected Location{309,350}, got Location{%f,%f}", g1.BlueTeam.Base.X, g1.BlueTeam.Base.Y)
	}
	if g1.BlueTeam.BaseRadius != 7.0 {
		t.Errorf("TestGenerateObjectives(4) failed, expected BaseRadius 7, got BaseRadius %f", g1.BlueTeam.BaseRadius)
	}
	if len(g1.ControlPoints) != 1 {
		t.Errorf("TestGenerateObjectives(5) failed, expected 1 ControlPoint, got %d", len(g1.ControlPoints))
	}
	if !inGameBounds(g1, g1.ControlPoints["CP1"].Location) {
		t.Errorf("TestGenerateObjectives(6) failed, expected inGameBounds(CP) to be TRUE, got FALSE")
	}
	if len(g1.Pickups) < 20 || len(g1.Pickups) > 25 {
		t.Errorf("TestGenerateObjectives(7) failed, expected between 20 and 25 pickups, got %d", len(g1.Pickups))
	}

	// a game arena that is a tall rectangle
	// MULTICAP mode with 2 control points
	g2 := &Game{
		Boundaries: []Border{
			{Location{0, 0}, Direction{100, 0}},
			{Location{100, 0}, Direction{0, 200}},
			{Location{100, 200}, Direction{-100, 0}},
			{Location{0, 200}, Direction{0, -200}},
		},
		RedTeam: &Team{Name:"RedTeam"},
		BlueTeam: &Team{Name:"BlueTeam"},
		Mode: MULTICAP,
	}
	g2.generateObjectives(2)
	time.Sleep(200 * time.Millisecond)
	if (g2.RedTeam.Base != Location{50.0, 191.0}) {
		t.Errorf("TestGenerateObjectives(8) failed, expected Location{50, 191}, got Location{%f,%f}", g2.RedTeam.Base.X, g2.RedTeam.Base.Y)
	}
	if g2.RedTeam.BaseRadius != 7.0 {
		t.Errorf("TestGenerateObjectives(9) failed, expected BaseRadius 7, got BaseRadius %f", g2.RedTeam.BaseRadius)
	}
	if (g2.BlueTeam.Base != Location{50.0, 9.0}) {
		t.Errorf("TestGenerateObjectives(10) failed, expected Location{50,9}, got Location{%f,%f}", g2.BlueTeam.Base.X, g2.BlueTeam.Base.Y)
	}
	if g2.BlueTeam.BaseRadius != 7.0 {
		t.Errorf("TestGenerateObjectives(11) failed, expected BaseRadius 7, got BaseRadius %f", g2.BlueTeam.BaseRadius)
	}
	if len(g2.ControlPoints) != 2 {
		t.Errorf("TestGenerateObjectives(12) failed, expected 2 ControlPoints, got %d", len(g2.ControlPoints))
	}
	if !inGameBounds(g2, g2.ControlPoints["CP1"].Location) {
		t.Errorf("TestGenerateObjectives(13) failed, expected inGameBounds(CP) to be TRUE, got FALSE")
	}
	if !inGameBounds(g2, g2.ControlPoints["CP2"].Location) {
		t.Errorf("TestGenerateObjectives(14) failed, expected inGameBounds(CP) to be TRUE, got FALSE")
	}
	if len(g2.Pickups) < 20 || len(g2.Pickups) > 25 {
		t.Errorf("TestGenerateObjectives(15) failed, expected between 20 and 25 pickups, got %d", len(g2.Pickups))
	}

	// a square arena
	// PAYLOAD mode
	g3 := &Game{
		Boundaries: []Border{
			{Location{0, 0}, Direction{100, 0}},
			{Location{100, 0}, Direction{0, 100}},
			{Location{100, 100}, Direction{-100, 0}},
			{Location{0, 100}, Direction{0, -100}},
		},
		RedTeam: &Team{Name:"RedTeam"},
		BlueTeam: &Team{Name:"BlueTeam"},
		Mode: PAYLOAD,
	}
	g3.generateObjectives(1)
	time.Sleep(200 * time.Millisecond)
	if (g3.RedTeam.Base != Location{50.0, 91.0}) {
		t.Errorf("TestGenerateObjectives(16) failed, expected Location{50, 91}, got Location{%f,%f}", g3.RedTeam.Base.X, g3.RedTeam.Base.Y)
	}
	if g3.RedTeam.BaseRadius != 7.0 {
		t.Errorf("TestGenerateObjectives(17) failed, expected BaseRadius 7, got BaseRadius %f", g3.RedTeam.BaseRadius)
	}
	if (g3.BlueTeam.Base != Location{50.0, 9.0}) {
		t.Errorf("TestGenerateObjectives(18) failed, expected Location{50,9}, got Location{%f,%f}", g3.BlueTeam.Base.X, g3.BlueTeam.Base.Y)
	}
	if g3.BlueTeam.BaseRadius != 7.0 {
		t.Errorf("TestGenerateObjectives(19) failed, expected BaseRadius 7, got BaseRadius %f", g3.BlueTeam.BaseRadius)
	}
	if len(g3.ControlPoints) != 1 {
		t.Errorf("TestGenerateObjectives(20) failed, expected 1 ControlPoint, got %d", len(g3.ControlPoints))
	}
	if (g3.ControlPoints["Payload"].Location != Location{50.0, 50.0}) {
		t.Errorf("TestGenerateObjectives(21) failed, expected Payload Location{50,50}, got Location{%f,%f}", g3.ControlPoints["Payload"].Location.X, g3.ControlPoints["Payload"].Location.Y)
	}
	if len(g3.Pickups) < 20 || len(g3.Pickups) > 25 {
		t.Errorf("TestGenerateObjectives(22) failed, expected between 20 and 25 pickups, got %d", len(g3.Pickups))
	}
	if (g3.PayloadPath != Direction{0.0, -82.0}) {
		t.Errorf("TestGenerateObjectives(23) failed, expected PayloadPath Direction{0,-86}, got Direction{%f,%f}", g3.PayloadPath.X, g3.PayloadPath.Y)
	}
	if g3.PayloadSpeed != 0.5 {
		t.Errorf("TestGenerateObjectives(24) failed, expected PayloadSpeed 0.5, got %f", g3.PayloadSpeed)
	}
}

func TestGeneratePickup(t *testing.T) {
	isTesting = true
	g := &Game{
		Boundaries: []Border{
			{Location{0, 0}, Direction{100, 0}},
			{Location{100, 0}, Direction{0, 100}},
			{Location{100, 100}, Direction{-100, 0}},
			{Location{0, 100}, Direction{0, -100}},
		},
		RedTeam: &Team{Name:"RedTeam"},
		BlueTeam: &Team{Name:"BlueTeam"},
	}
	generatePickup(g, 10.0, 10.0, 50.0, 10.0)
	if len(g.Pickups) != 1 {
		t.Errorf("TestGeneratePickup(1) failed, expected 1 pickup, got %d", len(g.Pickups))
	}
	if !inGameBounds(g, g.Pickups[0].Location) {
		t.Errorf("TestGeneratePickup(2) failed, expected inGameBounds(pickup) to be TRUE, got FALSE")
	}
	if g.Pickups[0].Location.X < 10.0 || g.Pickups[0].Location.X > 20.0 || g.Pickups[0].Location.Y < 10.0 || g.Pickups[0].Location.Y > 20.0 {
		t.Errorf("TestGeneratePickup(3) failed, expected X and Y-values in the [10,20] range, got (%f,%f)", g.Pickups[0].Location.X, g.Pickups[0].Location.Y)
	}
	generatePickup(g, 40.0, 60.0, 50.0, 10.0)
	if len(g.Pickups) != 2 {
		t.Errorf("TestGeneratePickup(4) failed, expected 2 pickups, got %d", len(g.Pickups))
	}
	if !inGameBounds(g, g.Pickups[1].Location) {
		t.Errorf("TestGeneratePickup(5) failed, expected inGameBounds(pickup) to be TRUE, got FALSE")
	}
	if g.Pickups[1].Location.X < 40.0 || g.Pickups[1].Location.X > 50.0 || g.Pickups[1].Location.Y < 60.0 || g.Pickups[1].Location.Y > 70.0 {
		t.Errorf("TestGeneratePickup(6) failed, expected X and Y-values in the [40,50] and [60,70] ranges, respectively, got (%f,%f)", g.Pickups[1].Location.X, g.Pickups[1].Location.Y)
	}
}

func TestNoIntersections(t *testing.T) {
	g := makeSampleGame()
	g.ControlPoints = nil
	g.ControlPoints = make(map[string]*ControlPoint)
	CP1 := ControlPoint{Location: Location{10.0,10.0}, Radius: 3.0}
	Pickup1 := PickupSpot{Location: Location{20.0,30.0}}
	g.ControlPoints["CP1"] = &CP1
	g.Pickups = append(g.Pickups, &Pickup1)
	loc1 := Location{14.0,10.0}
	loc2 := Location{20.0, 34.0}
	loc3 := Location{50.0,50.0}
	if noIntersections(g, loc1, 2.0) {
		t.Errorf("TestNoIntersections(1) failed, expected FALSE, got TRUE")
	}
	if !noIntersections(g, loc1, 0.5) {
		t.Errorf("TestNoIntersections(2) failed, expected TRUE, got FALSE")
	}
	if noIntersections(g, loc2, 2.0) {
		t.Errorf("TestNoIntersections(3) failed, expected FALSE, got TRUE")
	}
	if !noIntersections(g, loc2, 0.5) {
		t.Errorf("TestNoIntersections(4) failed, expected TRUE, got FALSE")
	}
	if !noIntersections(g, loc3, 5.0) {
		t.Errorf("TestNoIntersections(5) failed, expected TRUE, got FALSE")
	}
}
