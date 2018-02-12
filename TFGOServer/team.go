package main

type Team struct {
	Name string
	Players map[string]*Player
	Base Location
	Points int
}

func (t *Team) getPlayerLocs() []Location {
	return nil
}