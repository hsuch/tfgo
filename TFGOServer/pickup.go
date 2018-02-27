package main

import (
	"time"
	"math"
)

// consume a pickup and start its respawn timer
func (p *PickupSpot) consumePickup(player *Player, game *Game) {
	p.Pickup.use(player)
	p.Available = false
	sendPickupUpdate(p, game)
	p.SpawnTimer = time.AfterFunc(PICKUPRESPAWNTIME(), func() {
		p.Available = true
		p.SpawnTimer = nil
		sendPickupUpdate(p, game)
	})
}

// interface implementations
func (p ArmorPickup) use(player *Player) {
	player.Armor = intMin(MAXARMOR(), player.Armor + p.AP)
	sendVitalsUpdate(player)
}

func (p HealthPickup) use(player *Player) {
	player.Health = intMin(MAXHEALTH(), player.Health + p.HP)
	sendVitalsUpdate(player)
}

func (p WeaponPickup) use(player *Player) {
	sendAcquireWeapon(player, p.WP)
}

// functions for deciding on pickup attributes
func chooseArmorHealth(g *Game, loc Location, grange float64, isArmor bool) int {
	rLock.Lock()
	baseAH := r.Intn(2 * MAXARMOR())
	rLock.Unlock()
	var locAdj int
	if isArmor {
		locAdj = (int)(math.Floor(distance(g.findCenter(), loc) * (float64)(MAXARMOR())/grange))
	} else {
		locAdj = (int)(math.Floor((grange/2 - distance(g.findCenter(), loc)) * (float64)(MAXHEALTH())/grange))
	}
	armorHealth := intMax(baseAH-locAdj, 10)
	armorHealth = intMin(armorHealth, MAXARMOR())
	return armorHealth
}
