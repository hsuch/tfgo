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
	if testWeapon.canHit(Location{5, 5}, Location{4, 6}, Direction{-1, 1}) == math.MaxFloat64 {
		t.Errorf("TestCanHit(1) failed, expected distance < math.MaxFloat64 (Can hit), got distance == math.MaxFloat64 (Can't hit).")
	}

	// shot is right of target; within spread; within range
	if testWeapon.canHit(Location{5, 5}, Location{6, 6}, Direction{1, 1}) == math.MaxFloat64 {
		t.Errorf("TestCanHit(2) failed, expected distance < math.MaxFloat64 (Can hit), got distance == math.MaxFloat64 (Can't hit).")
	}

	// shot is left of target; outside spread; within range
	dist := testWeapon.canHit(Location{5, 5}, Location{4, 4}, Direction{-1, 1})
	if dist != math.MaxFloat64 {
		t.Errorf("TestCanHit(3) failed, expected Distance math.MaxFloat64 (Can't hit), got Distance %d (Can hit).", dist)
	}

	// shot is right of target; outside spread; within range
	dist = testWeapon.canHit(Location{5, 5}, Location{6, 4}, Direction{1, 1})
	if dist != math.MaxFloat64 {
		t.Errorf("TestCanHit(4) failed, expected Distance math.MaxFloat64 (Can't hit), got Distance %d (Can hit).", dist)
	}

	// shot is within spread; outside range
	dist = SWORD.canHit(Location{5, 5}, Location{10, 10}, Direction{1, 1})
	if dist != math.MaxFloat64 {
		t.Errorf("TestCanHit(5) failed, expected Distance math.MaxFloat64 (Can't hit), got Distance %d (Can hit).", dist)
	}

	// shot is outside spread; outside range
	dist = testWeapon.canHit(Location{5, 5}, Location{20, 20}, Direction{-1, -1})
	if dist != math.MaxFloat64 {
		t.Errorf("TestCanHit(6) failed, expected Distance math.MaxFloat64 (Can't hit), got Distance %d (Can hit).", dist)
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
			fname, pname, playerStatusToString[status], playerStatusToString[player.Status])
	}
}

func TestTakeHit(t *testing.T) {
	g := makeSampleGame()
	jenny := getJenny(g) // 100 hp, 0 armor
	oliver := getOliver(g) // 95 hp, 10 armor
	anders := getAnders(g) // 10 hp, 5 armor
	brad := getBrad(g) // 80 hp, 30 armor

	// no armor, hp > damage
	jenny.takeHit(g, SWORD)
	checkPlayerVitals(t, *jenny, jenny.Health - SWORD.Damage, 0, NORMAL, "TestTakeHit", "jenny")

	// armor < damage, hp + armor > damage
	oliver.takeHit(g, SWORD)
	checkPlayerVitals(t, *oliver, oliver.Health + oliver.Armor - SWORD.Damage, 0, NORMAL, "TestTakeHit", "oliver")

	// hp + armor < damage
	anders.takeHit(g, SWORD)
	checkPlayerVitals(t, *anders, 0, 0, RESPAWNING, "TestTakeHit", "anders")

	// armor > damage
	brad.takeHit(g, SWORD)
	checkPlayerVitals(t, *brad, brad.Health, brad.Armor - SWORD.Damage, NORMAL, "TestTakeHit", "brad")
}

func TestFire(t *testing.T) {
	g := makeSampleGame()
	jenny := getJenny(g)
	oliver := getOliver(g)
	anders := getAnders(g)
	brad := getBrad(g)

	// no targets hit
	brad.fire(g, SWORD,45)
	checkPlayerVitals(t, *jenny, jenny.Health, jenny.Armor, NORMAL, "TestFire", "brad->jenny")
	checkPlayerVitals(t, *oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "brad->oliver")
	checkPlayerVitals(t, *anders, anders.Health, anders.Armor, NORMAL, "TestFire", "brad->anders")

	// single target hit, non-fatal
	anders.fire(g, SWORD, 45)
	jenny.takeHit(g, SWORD)
	checkPlayerVitals(t, *jenny, jenny.Health, jenny.Armor, NORMAL, "TestFire", "anders->jenny")
	checkPlayerVitals(t, *oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "anders->oliver")
	checkPlayerVitals(t, *brad, brad.Health, brad.Armor, NORMAL, "TestFire", "anders->brad")

	// multiple targets available, hits closest
	oliver.handleLoc(g, Location{100, 99}, 0)
	anders.fire(g, SWORD, 45)
	jenny.takeHit(g, SWORD)
	checkPlayerVitals(t, *jenny, jenny.Health, jenny.Armor, NORMAL, "TestFire", "anders->jenny")
	checkPlayerVitals(t, *oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "anders->oliver")
	checkPlayerVitals(t, *brad, brad.Health, brad.Armor, NORMAL, "TestFire", "anders->brad")

	// single target hit, fatal
	jenny.fire(g, SWORD, 45)
	anders.takeHit(g, SWORD)
	checkPlayerVitals(t, *oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "jenny->oliver")
	checkPlayerVitals(t, *anders, anders.Health, anders.Armor, RESPAWNING, "TestFire", "jenny->anders")
	checkPlayerVitals(t, *brad, brad.Health, brad.Armor, NORMAL, "TestFire", "jenny->brad")
}
