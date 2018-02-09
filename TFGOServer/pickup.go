package TFGOServer

type Pickup interface {
	Use(game Game, player Player)
}
