package main

import "testing"

var directions = []Direction{
	{0, 1},
	{-1, 8},
	{1, 8},
	{-5, 0},
	{5, 0},
	{1, 12},
}

func TestGetPlayerLocs(t *testing.T) {
	// testing with a team with 2 members
	teamLocs := redTeam.getPlayerLocs()
	if teamLocs[0] != boundaries[2] {
		t.Errorf("TestGetTeamLocs(redTeam) failed, got: (%f,%f), want: (100,100).", teamLocs[0].X, teamLocs[0].Y)
	}
	if teamLocs[1] != boundaries[3] {
		t.Errorf("TestGetTeamLocs(redTeam) failed, got: (%f,%f), want: (0,100).", teamLocs[1].X, teamLocs[1].Y)
	}

	// testing with a team with no members
	teamLocs = Team{}.getPlayerLocs()
	if teamLocs != nil {
		t.Errorf("TestGetTeamLocs(emptyTeam) failed, expected output length 0, got length %d.", len(teamLocs))
	}
}

func TestTakeHit(t *testing.T) {
	/* HP, NO ARMOR, NO DEATH */
	p1 := jenny // 100 hp, 0 armor
	expected1 := jenny.Health - SWORD.Damage
	p1.takeHit (SWORD)
	if p1.Health !=  expected1 {
		t.Errorf("TestTakeHit(jenny) failed, got: (%d), want: (%d).", p1.Health, expected1)
	}
	if p1.Status != NORMAL {
		t.Errorf("TestTakeHit(jenny) failed, got (Status: %d), want: (Status: %d)", p1.Status, NORMAL)
	}

	/* HP, ARMOR < DMG, NO DEATH*/
	p2 := oliver // 104 hp, 2 armor
	p2.takeHit (SWORD)
	expected2_hp := oliver.Health + oliver.Armor - SWORD.Damage - oliver.Armor // splash
	expected2_armor := 0 // should go to 0 when less than damage

	// expect armor to go to 0 when dmg > armor
	if p2.Armor != expected2_armor  {
		t.Errorf("TestTakeHit(p2) failed, got (Armor: %d), want: (Armor: %d)", p2.Armor, expected2_armor)
	}
	if p2.Health != expected2_hp {
		t.Errorf("TestTakeHit(p2) failed, got (HP: %d), want: (HP: %d)", p2.Health, expected2_hp)
	}
	if p2.Status != NORMAL {
		t.Errorf("TestTakeHit(p2) failed, got (Status: %d), want: (Status: %d)", p2.Status, NORMAL)
	}

	/* HP, ARMOR < DMG, DEATH (Splash damage kills) */
	p3 := anders // 90 hp, 5 armor
	p3.Health = 20
	p3.takeHit (SWORD)
	if p3.Status != RESPAWNING {
		t.Errorf("TestTakeHit(anders) failed, got (Status: %d), want: (Status: %d)", p3.Status, RESPAWNING)
	}

	/* HP, ARMOR > DMG, NO DEATH */
	p4 := brad // 80 hp, 20 armor
	p4.Armor = 26
	p4.takeHit (SWORD)
	expected4_armor := p4.Armor - SWORD.Damage
	if !(p4.Armor > 0) {
		t.Errorf("TestTakeHit(p4) failed, got (Armor: %d), want: (Armor: %d)", p4.Armor, expected4_armor)
	}
}