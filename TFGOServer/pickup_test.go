package main

import "testing"

func TestConsumePickup(t *testing.T) {
	isTesting = true
	g := makeSampleGame()
	brad := getBrad(g)
	bradAP := brad.Armor // (HP 80, AP 30)
	ps := PickupSpot{Pickup: ArmorPickup{50}, Available: true}
	p := ps.Pickup
	ps.consumePickup(brad)
	if brad.Armor != bradAP + p.(ArmorPickup).AP {
		t.Errorf("TestPickupArmor(1) failed, expected (Armor: %d), got (Armor: %d)",
			bradAP + p.(*ArmorPickup).AP, brad.Armor)
	}
}

func TestPickupArmor(t *testing.T) {
	isTesting = true
	g := makeSampleGame()
	p := ArmorPickup{50}
	brad := getBrad(g) // (HP 80, AP 30)
	bradAP := brad.Armor
	p.use(brad)
	if brad.Armor != bradAP + p.AP {
		t.Errorf("TestPickupArmor(1) failed, expected (Armor: %d), got (Armor: %d)",
			bradAP + p.AP, brad.Armor)
	}
	p.use(brad)
	if brad.Armor != MAXARMOR() {
		t.Errorf("TestPickupArmor(2) failed, expected (Armor: %d), got (Armor: %d)",
			MAXARMOR(), brad.Armor)
	}
}

func TestPickupHealth(t *testing.T) {
	isTesting = true
	g := makeSampleGame()
	p := HealthPickup{50}
	anders := getAnders(g) // (HP 10, AP 5)
	andersHP := anders.Health
	p.use(anders)
	if anders.Health != andersHP + p.HP {
		t.Errorf("TestPickupHealth(1) failed, expected (Health: %d), got (Health: %d)",
			andersHP + p.HP, anders.Health)
	}
	p.use(anders)
	if anders.Health != MAXHEALTH() {
		t.Errorf("TestPickupHealth(2) failed, expected (Health: %d), got (Health: %d)",
			MAXHEALTH(), anders.Health)
	}
}
