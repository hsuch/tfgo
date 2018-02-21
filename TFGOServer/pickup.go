package main

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
func makeArmorPickup(loc Location) Pickup {
	return &ArmorPickup {50, loc}
}

func makeHealthPickup(loc Location) Pickup {
	return &HealthPickup {50, loc}
}

func makeWeaponPickup(wp Weapon, loc Location) Pickup {
	// we may want this to be random
	return &WeaponPickup {wp, loc}
}