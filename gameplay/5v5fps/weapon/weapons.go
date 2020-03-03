package weapon

type Weapon struct {
	ID int
	Name string
	Damage uint
	Price uint
	WT WeaponType
}

type WeaponType int

const (
	Primary WeaponType = 0
	Secondary WeaponType = 1
	Knife WeaponType = 2
)

// WeaponsFromFile reads a json file and returns a slice of Weapon
func WeaponsFromFile(path string) ([]Weapon, error) {
	var weapons []Weapon
	
	contents, err := ioutil.ReadFile()
	if err != nil {
		return weapons, err
	}
	err := json.Unmarshal(contents, &weapons)
	return weapons, err
}