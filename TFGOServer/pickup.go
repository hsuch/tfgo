package main

type Pickup interface {
	use(game *Game, player *Player)
}
