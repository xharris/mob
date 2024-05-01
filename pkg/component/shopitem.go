package component

type ShopItem struct {
	AddMods     []Mod
	RemoveMods  []Mod
	UpgradeMods []Mod
	Cost        int
	Name        string
	Purchased   bool
}
