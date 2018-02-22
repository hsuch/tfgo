package main

import "testing"

func TestPickupArmor(t *testing.T) {
	isTesting = true
	g := makeSampleGame()
	p := makeArmorPickup(50);
	brad := getBrad(g) // (HP 80, AP 30)
	bradAP := brad.Armor
	p.use(brad)
	if brad.Armor != bradAP + p.(*ArmorPickup).AP {
		t.Errorf("TestPickupArmor(1) failed, expected (Armor: %d), got (Armor: %d)",
			bradAP + p.(*ArmorPickup).AP, brad.Armor)
	}
	p.use(brad)
	if brad.Armor != MAXARMOR() {
		t.Errorf("%TestPickupArmor(2) failed, expected (Armor: %d), got (Armor: %d)",
			MAXARMOR(), brad.Armor)
	}
}

func TestPickupHealth(t *testing.T) {
	isTesting = true
	g := makeSampleGame()
	p := makeHealthPickup(50)
	anders := getAnders(g) // (HP 10, AP 5)
	andersHP := anders.Health
	p.use(anders)
	if anders.Health != andersHP + p.(*HealthPickup).HP {
		t.Errorf("TestPickupHealth(1) failed, expected (Health: %d), got (Health: %d)",
			andersHP + p.(*HealthPickup).HP, anders.Health)
	}
	p.use(anders)
	if anders.Health != MAXHEALTH() {
		t.Errorf("%TestPickupHealth(2) failed, expected (Health: %d), got (Health: %d)",
			MAXHEALTH(), anders.Health)
	}
}

func TestPickupWeapon(t *testing.T) {
	isTesting = true
	g := makeSampleGame()
	p := makeWeaponPickup(SHOTGUN);
	anders := getAnders(g) // (HP 10, AP 5)
	if _, ok := anders.Weapons["Shotgun"]; ok {
		t.Errorf("TestPickupWeapon(1) failed, expected Shotgun to NOT be in Weapons list, got otherwise")
	}
	p.use(anders)
	if _, ok := anders.Weapons["Shotgun"]; !ok {
		t.Errorf("TestPickupWeapon(2) failed, expected Shotgun to be in Weapons list, got otherwise")
	}
}
