package main

import (
	"time"
	"math"
)

// consume a pickup, set status to respawning, call use () on it
func (p *PickupSpot) consumePickup(player *Player) {
	if (p.Available) {
		p.Pickup.use(player)
		go p.awaitRespawn()
	}
}

func (p *PickupSpot) awaitRespawn() {
	p.Available = false
	if (!isTesting) {
		p.SpawnTimer = time.AfterFunc(PICKUPRESPAWNTIME(),
			func() {
				p.respawn()
			})
	}
}

func (p *PickupSpot) respawn() {
	p.Available = true
	p.SpawnTimer = nil
}

// Interface implementations
func (p ArmorPickup) use(player *Player) {
	player.Armor = intMin(MAXARMOR(), player.Armor + p.AP)
}

func (p HealthPickup) use(player *Player) {
	player.Health = intMin(MAXHEALTH(), player.Health + p.HP)
}

func (p WeaponPickup) use(player *Player) {
	player.Weapons[p.WP.Name] = p.WP
}


// Functions for spawning pickups
func makeArmorPickup(val int) Pickup {
	return &ArmorPickup {val}
}

func makeHealthPickup(val int) Pickup {
	return &HealthPickup {val}
}

func makeWeaponPickup(wp Weapon) Pickup {
	// we may want this to be random
	return &WeaponPickup {wp}
}

//Functions for deciding on pickup attributes
func chooseArmorHealth(g *Game, loc Location, grange float64) int {
	rLock.Lock()
	base_ah := r.Intn(2 * MAXARMOR())
	rLock.Unlock()
	loc_adj := (int)(math.Floor(distance(g.findCenter(), loc) * (float64)(MAXARMOR())/grange))
	armor_health := intMax(base_ah - loc_adj, 10)
	armor_health = intMin(armor_health, MAXARMOR())
	return armor_health
}
