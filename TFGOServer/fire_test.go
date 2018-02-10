package main

import (
	"testing"
	"math"
)

var testWeapon = Weapon{
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
	// red team (two members)
	playerLocs := redTeam.getPlayerLocs()
	if playerLocs[0] != boundaries[2] {
		t.Errorf("TestGetTeamLocs(redTeam) failed, expected (100,100), got (%f,%f).", playerLocs[0].X, playerLocs[0].Y)
	}
	if playerLocs[1] != boundaries[3] {
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
			fname, pname, playerStatusString(status), playerStatusString(player.Status))
	}
}

func TestTakeHit(t *testing.T) {
	// no armor, hp > damage
	p1 := jenny // 100 hp, 0 armor
	p1.takeHit(SWORD)
	checkPlayerVitals(t, p1, jenny.Health - SWORD.Damage, 0, NORMAL, "TestTakeHit", "jenny")

	// armor < damage, hp + armor > damage
	p2 := oliver // 104 hp, 2 armor
	p2.takeHit (SWORD)
	checkPlayerVitals(t, p2, oliver.Health + oliver.Armor - SWORD.Damage, 0, NORMAL, "TestTakeHit", "oliver")

	// hp + armor < damage
	p3 := anders // 10 hp, 5 armor
	p3.takeHit(SWORD)
	checkPlayerVitals(t, p3, 0, 0, RESPAWNING, "TestTakeHit", "anders")

	// armor > damage
	p4 := brad // 80 hp, 30 armor
	p4.takeHit (SWORD)
	checkPlayerVitals(t, p4, brad.Health, brad.Armor - SWORD.Damage, NORMAL, "TestTakeHit", "brad")
}

func TestFire(t *testing.T) {
	p1 := jenny
	p2 := oliver
	p3 := anders
	p4 := brad

	// no targets hit
	p4.fire(Direction{1, 1})
	checkPlayerVitals(t, p1, jenny.Health, jenny.Armor, NORMAL, "TestFire", "brad->jenny")
	checkPlayerVitals(t, p2, oliver.Health, oliver.Armor, NORMAL, "TestFire", "brad->oliver")
	checkPlayerVitals(t, p3, anders.Health, anders.Armor, NORMAL, "TestFire", "brad->anders")

	// single target hit, non-fatal
	p3.fire(Direction{1, 1})
	p1check := jenny
	p1check.takeHit(SWORD)
	checkPlayerVitals(t, p1, p1check.Health, p1check.Armor, NORMAL, "TestFire", "anders->jenny")
	checkPlayerVitals(t, p2, oliver.Health, oliver.Armor, NORMAL, "TestFire", "anders->oliver")
	checkPlayerVitals(t, p4, brad.Health, brad.Armor, NORMAL, "TestFire", "anders->brad")

	// multiple targets available, hits closest
	p2move := oliver
	p2move.handleLoc(Location{100, 99})
	p3.fire(Direction{1, 1})
	p1check.takeHit(SWORD)
	checkPlayerVitals(t, p1, p1check.Health, p1check.Armor, NORMAL, "TestFire", "anders->jenny")
	checkPlayerVitals(t, p2move, oliver.Health, oliver.Armor, NORMAL, "TestFire", "anders->oliver")
	checkPlayerVitals(t, p4, brad.Health, brad.Armor, NORMAL, "TestFire", "anders->brad")

	// single target hit, fatal
	p1.fire(Direction{1, 1})
	p3check := anders
	p3check.takeHit(SWORD)
	checkPlayerVitals(t, p2, oliver.Health, oliver.Armor, NORMAL, "TestFire", "jenny->oliver")
	checkPlayerVitals(t, p3, p3check.Health, p3check.Armor, RESPAWNING, "TestFire", "jenny->anders")
	checkPlayerVitals(t, p4, brad.Health, brad.Armor, NORMAL, "TestFire", "jenny->brad")
}
