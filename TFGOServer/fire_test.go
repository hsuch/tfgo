package main

// fire_test.go: tests for weapon firing related functions in fire.go

import (
	"testing"
	"math"
)

var testWeapon = Weapon {
	Spread: math.Pi/2,
	Range: 5,
}

/* checkPlayerVitals - helper function that checks a player's health, armor,
 *                     and status to make sure the appropriate values are set
 */
func checkPlayerVitals(t *testing.T, player *Player, hp, armor int, status PlayerStatus, fname, pname string) {
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

func TestDot(t *testing.T) {
	v := Direction{X : float64(-1), Y : float64(2)}
	w := Direction{X : float64(-3), Y : float64(4)}
	expect := float64 (11)
	got := dot(v, w)
	if !isAcceptableError (got, expect, 0) {
		t.Errorf ("TestDot failed, expected %g, got %g", expect, got)
	}
}

func TestMagnitude(t *testing.T) {
	v := Direction{X : float64(1), Y : float64(2)}
	got := v.magnitude ()
	expect := math.Sqrt(float64(5))
	if !isAcceptableError(got, expect, 0) {
		t.Errorf ("TestMagnitude failed, expected %g, got %g", expect, got)
	}
}

func TestCanHit(t *testing.T) {
	// shot is left of target; within spread; within range
	dist := testWeapon.canHit(Location{5, 5}, Location{4.5, 6}, Direction{-1, 1})
	if dist == math.MaxFloat64 {
		t.Errorf("TestCanHit(1) failed, expected CAN HIT, got CAN'T HIT. (Dist: math.MaxFloat64)")
	}

	// shot is right of target; within spread; within range
	dist = testWeapon.canHit(Location{5, 5}, Location{5.5, 6}, Direction{1, 1})
	if dist == math.MaxFloat64 {
		t.Errorf("TestCanHit(2) failed, expected CAN HIT, got CAN'T HIT. (Dist: math.MaxFloat64)")
	}

	// shot is left of target; outside spread; within range
	dist = testWeapon.canHit(Location{5, 5}, Location{5.5, 6}, Direction{-1, 1})
	if dist != math.MaxFloat64 {
		t.Errorf("TestCanHit(3) failed, expected CAN'T HIT, got CAN HIT. (Dist: %d)", dist)
	}

	// shot is right of target; outside spread; within range
	dist = testWeapon.canHit(Location{5, 5}, Location{4.5, 6}, Direction{1, 1})
	if dist != math.MaxFloat64 {
		t.Errorf("TestCanHit(4) failed, expected CAN'T HIT, got CAN HIT. (Dist: %d)", dist)
	}

	// shot is within spread; outside range
	dist = testWeapon.canHit(Location{5, 5}, Location{5, 11}, Direction{0, 1})
	if dist != math.MaxFloat64 {
		t.Errorf("TestCanHit(5) failed, expected CAN'T HIT, got CAN HIT. (Dist: %d)", dist)
	}

	// shot is outside spread; outside range
	dist = testWeapon.canHit(Location{5, 5}, Location{20, 20}, Direction{0, -1})
	if dist != math.MaxFloat64 {
		t.Errorf("TestCanHit(6) failed, expected CAN'T HIT, got CAN HIT. (Dist: %d)", dist)
	}
}

func TestFire(t *testing.T) {
	g := makeSampleGame()
	jenny := getJenny(g) // (49, 75)
	oliver := getOliver(g) // (50, 75)
	anders := getAnders(g) // (49.5, 75)
	brad := getBrad(g) // (49, 24)

	// no targets hit
	brad.fire(g, SWORD,45)
	checkPlayerVitals(t, jenny, jenny.Health, jenny.Armor, NORMAL, "TestFire", "brad->jenny")
	checkPlayerVitals(t, oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "brad->oliver")
	checkPlayerVitals(t, anders, anders.Health, anders.Armor, NORMAL, "TestFire", "brad->anders")

	// single target hit, non-fatal
	anders.fire(g, SWORD, 45)
	jenny.takeHit(g, SWORD)
	checkPlayerVitals(t, jenny, jenny.Health, jenny.Armor, NORMAL, "TestFire", "anders->jenny")
	checkPlayerVitals(t, oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "anders->oliver")
	checkPlayerVitals(t, brad, brad.Health, brad.Armor, NORMAL, "TestFire", "anders->brad")

	// multiple targets available, hits closest
	oliver.updateLocation(g, Location{100, 99}, 0)
	anders.fire(g, SWORD, 45)
	jenny.takeHit(g, SWORD)
	checkPlayerVitals(t, jenny, jenny.Health, jenny.Armor, NORMAL, "TestFire", "anders->jenny")
	checkPlayerVitals(t, oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "anders->oliver")
	checkPlayerVitals(t, brad, brad.Health, brad.Armor, NORMAL, "TestFire", "anders->brad")

	// single target hit, fatal
	jenny.fire(g, SWORD, 45)
	anders.takeHit(g, SWORD)
	checkPlayerVitals(t, oliver, oliver.Health, oliver.Armor, NORMAL, "TestFire", "jenny->oliver")
	checkPlayerVitals(t, anders, anders.Health, anders.Armor, RESPAWNING, "TestFire", "jenny->anders")
	checkPlayerVitals(t, brad, brad.Health, brad.Armor, NORMAL, "TestFire", "jenny->brad")
}

func TestTakeHit(t *testing.T) {
	g := makeSampleGame()
	jenny := getJenny(g) // 100 hp, 0 armor
	oliver := getOliver(g) // 95 hp, 10 armor
	anders := getAnders(g) // 10 hp, 5 armor
	brad := getBrad(g) // 80 hp, 30 armor

	// no armor, hp > damage
	jenny.takeHit(g, SWORD)
	checkPlayerVitals(t, jenny, jenny.Health - SWORD.Damage, 0, NORMAL, "TestTakeHit", "jenny")

	// armor < damage, hp + armor > damage
	oliver.takeHit(g, SWORD)
	checkPlayerVitals(t, oliver, oliver.Health + oliver.Armor - SWORD.Damage, 0, NORMAL, "TestTakeHit", "oliver")

	// hp + armor < damage
	anders.takeHit(g, SWORD)
	checkPlayerVitals(t, anders, 0, 0, RESPAWNING, "TestTakeHit", "anders")

	// armor > damage
	brad.takeHit(g, SWORD)
	checkPlayerVitals(t, brad, brad.Health, brad.Armor - SWORD.Damage, NORMAL, "TestTakeHit", "brad")
}

func TestAwaitRespawn(t *testing.T) {
	// func (p *Player) awaitRespawn(game *Game)
}

func TestRespawn(t *testing.T) {
	// func (p *Player) respawn(game *Game)
}
