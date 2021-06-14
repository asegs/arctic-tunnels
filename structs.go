package main

type BodyArmor int
type ItemType int


const (
	HEAD BodyArmor = iota
	BODY
	HANDS
	LEGS
	FEET
)

const (
	WEAPON ItemType = iota
	CLOTHING
	EQUIPMENT
	QUEST
)


type Character struct {
	Name string `json:"name"`
	DefaultHealth int `json:"default_health"`
	Health int `json:"health"`
	Hunger int `json:"hunger"`
	Thirst int `json:"thirst"`
	Rubles int `json:"rubles"`
	Inventory []InventoryItem



}


type InventoryItem struct {
	Item Sellable `json:"item"`
	Type ItemType `json:"type"`
}

type Sellable interface {
	evaluate()int
	sell(seller *Character,buyer *Character)
}

type Modifier struct {
	StatEffected string `json:"stat_effected"`
	Modifier float64 `json:"modifier"`
	IsPercent bool `json:"is_percent"`
	Duration int `json:"duration"`
}


type Clothing struct {
	Name string `json:"name"`
	ColdProtection int `json:"cold_protection"`
	HeatProtection int `json:"heat_protection"`
	Bulletproof int `json:"bulletproof"`
	MeleeProtection int `json:"melee_protection"`
	Weight int `json:"weight"`
	Durability int `json:"durability"`
	Modifier Modifier `json:"modifier"`
}

type Weapon interface {
	attack(attacker *Character, defender *Character)
}


type Gun struct{
	Name string `json:"name"`
	Calibre string `json:"calibre"`
	EffectiveRange string `json:"effective_range"`
	MaxDamage int `json:"damage"`
	Accuracy int `json:"accuracy"`
	LoadedMagazine Ammo `json:"loaded_magazine"`
	Weight int `json:"weight"`
	Durability int `json:"durability"`
}

type Ammo struct {
	Name string `json:"name"`
	Calibre string `json:"calibre"`
	ArmorPiercing int `json:"armor_piercing"`
}

type Melee struct{
	Name string `json:"name"`
	Range int `json:"range"`
	PiercingDamage int `json:"piercing_damage"`
	SlashingDamage int `json:"slashing_damage"`
	CrushingDamage int `json:"crushing_damage"`
	Weight int `json:"weight"`
	Durability int `json:"durability"`
}

type Coordinate struct{
	Row int `json:"row"`
	Col int `json:"col"`
}

type Place struct {

}

func (g Gun)evaluate()int{
	//really do calculations
	return 100
}

func (c Clothing)evaluate()int{
	//really do calculations
	return 100
}

func (g Gun)sell(seller *Character,buyer *Character){
	price := g.evaluate()
	seller.Rubles+=price
	buyer.Rubles-=price
	//seller.remove(item)
	buyer.Inventory = append(buyer.Inventory,InventoryItem{
		Item: g,
		Type: WEAPON,
	})

}

func (c Clothing)sell(seller *Character,buyer *Character){
	price := c.evaluate()
	seller.Rubles+=price
	buyer.Rubles-=price
	//seller.remove(item)
	buyer.Inventory = append(buyer.Inventory,InventoryItem{
		Item: c,
		Type: CLOTHING,
	})
}


