package main

import (
	"time"
	)

// consume a pickup, set status to respawning, call use () on it
func (p *PickupSpot) consumePickup (game *Game, player *Player) {
	// if (!p.Available) // pickup not active / live, already used
		// return

	// }
	// else {
	if (p.Available) {
		p.Pickup.use (player)
		go p.awaitRespawn (game)
	}
}

func (p *PickupSpot) awaitRespawn (game *Game) {
	p.Available = false
	if (!isTesting) {
		p.SpawnTimer = time.AfterFunc(PICKUPRESPAWNTIME (),
			func() {p.Available = true })
		p.respawn (game)
	}
}

func (p *PickupSpot) respawn (game *Game) {
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
	loc_adj := (int)math.floor(distance(g.findCenter(), loc) * MAXARMOR()/grange)
	armor_health := intMax(base_a - loc_adj, 10)
	armor_health = intMin(armor_health, MAXARMOR())
	return armor_health
}
