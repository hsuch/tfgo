package main

import (
	"time"
	"math"
)

// consume a pickup and start its respawn timer
func (p *PickupSpot) consumePickup(player *Player) {
	p.Pickup.use(player)
	p.Available = false
	p.SpawnTimer = time.AfterFunc(PICKUPRESPAWNTIME(), func() {
		p.Available = true
		p.SpawnTimer = nil
	})
}

// interface implementations
func (p ArmorPickup) use(player *Player) {
	player.Armor = intMin(MAXARMOR(), player.Armor + p.AP)
	sendVitalStats(player)
}

func (p HealthPickup) use(player *Player) {
	player.Health = intMin(MAXHEALTH(), player.Health + p.HP)
	sendVitalStats(player)
}

func (p WeaponPickup) use(player *Player) {
	sendWeaponAcquire(player, p.WP)
}

// functions for deciding on pickup attributes
func chooseArmorHealth(g *Game, loc Location, grange float64) int {
	rLock.Lock()
	baseAH := r.Intn(2 * MAXARMOR())
	rLock.Unlock()
	locAdj := (int)(math.Floor(distance(g.findCenter(), loc) * (float64)(MAXARMOR())/grange))
	armorHealth := intMax(baseAH-locAdj, 10)
	armorHealth = intMin(armorHealth, MAXARMOR())
	return armorHealth
}
