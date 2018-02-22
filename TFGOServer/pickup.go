package main

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
func makeArmorPickup() Pickup {
	return &ArmorPickup {50}
}

func makeHealthPickup() Pickup {
	return &HealthPickup {50}
}

func makeWeaponPickup(wp Weapon) Pickup {
	// we may want this to be random
	return &WeaponPickup {wp}
}