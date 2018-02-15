package main

// setup_test.go: tests for game setup functions in setup.go

import "testing"

func TestFindCenter(t *testing.T) {
	// func (g *Game) findCenter() Location
}

func TestCreatePlayer(t *testing.T) {
	// func createPlayer(conn net.Conn, name, icon string) *Player
}

func TestCreateGameID(t *testing.T) {
	// func createGameID() string
}

func TestSetBoundaries(t *testing.T) {
	// func (g *Game) setBoundaries(boundaries []interface{})
}

func TestCreateGame(t *testing.T) {
	// func createGame(conn net.Conn, data map[string]interface{}) (*Game, *Player)
}

func TestJoinGame(t *testing.T) {
	// func (p *Player) joinGame(gameID string) *Game
}

func TestGenerateObjectives(t *testing.T) {
	// func (g *Game) generateObjectives(numCP int)
}

func TestRandomizeTeams(t *testing.T) {
	// func (g *Game) randomizeTeams()
}

func TestStart(t *testing.T) {
	// func (g *Game) start()
}

func TestStop(t *testing.T) {
	// func (g *Game) stop()
}
