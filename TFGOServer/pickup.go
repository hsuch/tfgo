package TFGOServer

type Pickup interface {
	use(game *Game, player *Player)
}
