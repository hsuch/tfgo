package main

import "testing"

func TestPickupArmor(t *testing.T) {
	isTesting = true
	g := makeSampleGame()
	p := makeArmorPickup(Location{1, 1});
	brad := getBrad(g) // (HP 80, AP 30)
	bradAP := brad.Armor
	p.use(g, brad)
	if brad.Armor != bradAP + p.(*ArmorPickup).AP {
		t.Errorf("TestPickupArmor(1) failed, expected (Armor: %d), got (Armor: %d)",
			bradAP + p.(*ArmorPickup).AP, brad.Armor)
	}
	p.use(g, brad)
	if brad.Armor != 100 {
		t.Errorf("%TestPickupArmor(2) failed, expected (Armor: %d), got (Armor: %d)",
			100, brad.Armor)
	}
}

