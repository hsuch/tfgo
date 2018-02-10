package TFGOServer

type Team struct {
	Name string
	Players []*Player
	Base Location
	Points int
}

func (t *Team) getPlayerLocs() []Location {
	return nil
}