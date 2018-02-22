package main

// consume a pickup, set status to respawning, call use () on it
func (p *PickupSpot) consumePickup (game *Game, player *Player) {
	if (p.Status == pRESPAWNING) {
		return
	}
	p.Pickup.use (game, player)
	go p.awaitRespawn (game)
}

func (p *PickupSpot) awaitRespawn (game *Game) {
	p.Status = pRESPAWNING
	if (!isTesting) {
		<- p.StatusTimer.C
		p.respawn (game)
	}
	// send status update? 
}

func (p *PickupSpot) respawn (game *Game) {
	p.Status = pNORMAL
	p.StatusTimer = nil
	// send status update? 
}

// Interface implementations
func (p ArmorPickup) use(game *Game, player *Player) {
	player.Armor = intMin(MAXARMOR(), player.Armor + p.AP)
}

func (p HealthPickup) use(game *Game, player *Player) {
	player.Health = intMin(MAXHEALTH(), player.Health + p.HP)
}

func (p WeaponPickup) use(game *Game, player *Player) {
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