package main

// setup_test.go: tests for game setup functions in setup.go

import "testing"

func TestFindCenter(t *testing.T) {
	g := makeSampleGame()
	x := g.findCenter().X
	y := g.findCenter().Y
	if !isAcceptableError(x, 50, EPSILON) && isAcceptableError(y, 50, EPSILON) {
		t.Errorf("TestFindCenter(1) failed, expected Location (50, 50), got Location (%d, %d)", x, y)
	}
}

func TestCreatePlayer(t *testing.T) {
	p := createPlayer(nil, "Alice", "testIcon")
	checkPlayerVitals(t, p, 100, 0, NORMAL, "TestCreatePlayer", "Alice")
	if p.Icon != "testIcon" {
		t.Errorf("TestCreatePlayer(1) failed, expected Icon testIcon, got Icon %s", p.Icon)
	}
}

func TestSetBoundaries(t *testing.T) {
	// func (g *Game) setBoundaries(boundaries []interface{})
}

func TestCreateGame(t *testing.T) {
	// func createGame(conn net.Conn, data map[string]interface{}) (*Game, *Player)

	//map[string]interface{} {
	//	"Name" : "Game 1",
	//	"Password": "Game1Pass",
	//}
}

func TestGenerateObjectives(t *testing.T) {
	// func (g *Game) generateObjectives(numCP int)
}
