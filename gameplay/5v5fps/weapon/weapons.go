package weapon

import (
	"encoding/json"
	"io/ioutil"
)

type Weapon struct {
	ID      int
	Name    string
	Damage  uint
	Price   uint
	WT      WeaponType
	AmmoMax int
	Ammo    int
}

type WeaponType int

const (
	Primary   WeaponType = 0
	Secondary WeaponType = 1
	Knife     WeaponType = 2
)

var (
	SecondaryOne = Weapon{1, "secondary_one", 30, 200, Secondary, 20, 20}
	PrimaryOne = Weapon{2, "primary_one", 50, 1000, Primary, 30, 30}
	KnifeOne = Weapon{3, "knife_one", 40, 300, Knife, 0, 0}
)

// WeaponsFromFile reads a json file and returns a slice of Weapon
func WeaponsFromFile(path string) ([]Weapon, error) {
	var weapons []Weapon

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return weapons, err
	}
	err = json.Unmarshal(contents, &weapons)
	return weapons, err
}
