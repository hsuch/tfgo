package main

import (
	"testing"
	"math"
)

var testWeapon = Weapon {
	Spread: math.Pi/2,
	Range: 5,
}

func TestCanHit(t *testing.T) {
	// shot is left of target; within spread; within range
	if !testWeapon.canHit(Location{5, 5}, Location{4, 6}, Direction{-1, 1}) {
		t.Errorf("TestCanHit(1) failed, expected TRUE, got FALSE.")
	}

	// shot is right of target; within spread; within range
	if !testWeapon.canHit(Location{5, 5}, Location{6, 6}, Direction{1, 1}) {
		t.Errorf("TestCanHit(2) failed, expected TRUE, got FALSE.")
	}

	// shot is left of target; outside spread; within range
	if testWeapon.canHit(Location{5, 5}, Location{4, 4}, Direction{-1, 1}) {
		t.Errorf("TestCanHit(3) failed, expected FALSE, got TRUE.")
	}

	// shot is right of target; outside spread; within range
	if testWeapon.canHit(Location{5, 5}, Location{6, 4}, Direction{1, 1}) {
		t.Errorf("TestCanHit(4) failed, expected FALSE, got TRUE.")
	}

	// shot is within spread; outside range
	if SWORD.canHit(Location{5, 5}, Location{10, 10}, Direction{1, 1}) {
		t.Errorf("TestCanHit(5) failed, expected FALSE, got TRUE.")
	}

	// shot is outside spread; outside range
	if testWeapon.canHit(Location{5, 5}, Location{20, 20}, Direction{-1, -1}) {
		t.Errorf("TestCanHit(6) failed, expected FALSE, got TRUE.")
	}
}

func TestGetPlayerLocs(t *testing.T) {
	g := makeSampleGame()
	redTeam := g.RedTeam

	// red team (two members)
	playerLocs := redTeam.getPlayerLocs()
	if playerLocs[0] != g.Boundaries[2] {
		t.Errorf("TestGetTeamLocs(redTeam) failed, expected (100,100), got (%f,%f).", playerLocs[0].X, playerLocs[0].Y)
	}
	if playerLocs[1] != g.Boundaries[3] {
		t.Errorf("TestGetTeamLocs(redTeam) failed, expected (0,100), got (%f,%f).", playerLocs[1].X, playerLocs[1].Y)
	}

	// empty team (no members)
	playerLocs = Team{}.getPlayerLocs()
	if playerLocs != nil {
		t.Errorf("TestGetTeamLocs(emptyTeam) failed, expected output length 0, got length %d.", len(playerLocs))
	}
}

func checkPlayerVitals(t *testing.T, player Player, hp, armor int, status PlayerStatus, fname, pname string) {
	if player.Health != hp {
		t.Errorf("%s(%s) failed, expected (Health: %d), got (Health: %d)",
			fname, pname, hp, player.Health)
	}

	if player.Armor != armor {
		t.Errorf("%s(%s) failed, expected (Armor: %d), got (Armor: %d)",
			fname, pname, armor, player.Armor)
	}

	if player.Status != status {
		t.Errorf("%s(%s) failed, expected (Status: %s), got (Status: %s)",
			fname, pname, PlayerStatusMap[status], PlayerStatusMap[player.Status])
	}
}

func TestTakeHit(t *testing.T) {
	g := makeSampleGame()
	jenny := getJenny(g) // 100 hp, 0 armor
	oliver := getOliver(g) // 95 hp, 10 armor
	anders := getAnders(g) // 10 hp, 5 armor
	brad := getBrad(g) // 80 hp, 30 armor

	// no armor, hp > damage
	jenny.takeHit(SWORD)
	checkPlayerVitals(t, *jenny, jenny.Health - SWORD.Damage, 0, NORMAL, "TestTakeHit", "jenny")

	// armor < damage, hp + armor > damage
	oliver.takeHit (SWORD)
	checkPlayerVitals(t, *oliver, oliver.Health + oliver.Armor - SWORD.Damage, 0, NORMAL, "TestTakeHit", "oliver")

	// hp + armor < damage
	anders.takeHit(SWORD)
	checkPlayerVitals(t, *anders, 0, 0, RESPAWNING, "TestTakeHit", "anders")

	// armor > damage
	brad.takeHit (SWORD)
	checkPlayerVitals(t, *brad, brad.Health, brad.Armor - SWORD.Damage, NORMAL, "TestTakeHit", "brad")
}

func TestFire(t *testing.T) {
	g := makeSampleGame()
	jenny := getJenny(g)
	oliver := getOliver(g)
	anders := getAnders(g)
	brad := getBrad(g)

	// no targets hit
	brad.fire(g, Direction{1, 1})
	checkPlayerVitals(t, *jenny, jenny.Health, jenny.Armor, NORMAL, "TestFire", "brad->jenny")
	checkPlayerVitals(t, *oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "brad->oliver")
	checkPlayerVitals(t, *anders, anders.Health, anders.Armor, NORMAL, "TestFire", "brad->anders")

	// single target hit, non-fatal
	anders.fire(g, Direction{1, 1})
	jenny.takeHit(SWORD)
	checkPlayerVitals(t, *jenny, jenny.Health, jenny.Armor, NORMAL, "TestFire", "anders->jenny")
	checkPlayerVitals(t, *oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "anders->oliver")
	checkPlayerVitals(t, *brad, brad.Health, brad.Armor, NORMAL, "TestFire", "anders->brad")

	// multiple targets available, hits closest
	oliver.handleLoc(g, Location{100, 99})
	anders.fire(g, Direction{1, 1})
	jenny.takeHit(SWORD)
	checkPlayerVitals(t, *jenny, jenny.Health, jenny.Armor, NORMAL, "TestFire", "anders->jenny")
	checkPlayerVitals(t, *oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "anders->oliver")
	checkPlayerVitals(t, *brad, brad.Health, brad.Armor, NORMAL, "TestFire", "anders->brad")

	// single target hit, fatal
	jenny.fire(g, Direction{1, 1})
	anders.takeHit(SWORD)
	checkPlayerVitals(t, *oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "jenny->oliver")
	checkPlayerVitals(t, *anders, anders.Health, anders.Armor, RESPAWNING, "TestFire", "jenny->anders")
	checkPlayerVitals(t, *brad, brad.Health, brad.Armor, NORMAL, "TestFire", "jenny->brad")
}
