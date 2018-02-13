package main

type Team struct {
	Name string
	Players map[string]*Player
	Base Location
	Points int
}

func (t *Team) getPlayerLocs() []Location {
	var locs []Location
	for _, p := range t.Players {
		locs = append(locs, p.Location)
	}
	return locs
}
