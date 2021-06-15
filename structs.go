package main

import "fmt"

type LogMode int

type BodyArmor int
type ItemType int
type CharType int


const MAX_GUN_WEIGHT = 24
const MAX_GUN_RELOAD_TIME = 4000
const MAX_GUN_TIME_BETWEEN_SHOTS = 1500
const BULLET_MISS_HARSHNESS = 1.75

const (
	OFF LogMode = iota
	BASIC
	DEBUG
)

const (
	HEAD BodyArmor = iota
	CHEST
	HANDS
	LEGS
	FEET
)

var BODY_SCORES = [...]float64{1.5,1.0,0.25,0.35,0.45}
var BODY_MAX_WEIGHTS = [...]float64{8,40,1.5,6,5}
var BODY_HIT_MODIFIERS = [...]float64{0.5,2.0,0.65,1.25,1.0}

const (
	WEAPON ItemType = iota
	CLOTHING
	EQUIPMENT
	QUEST
)

const (
	PLAYER CharType = iota
	FRIENDLY
	LONER
	MILITARY
	LOOTER
	ANIMAL
	MONSTER
)


const LOG_MODE = LogMode(DEBUG)

const RUBLE = '₽'


type Character struct {
	Name string `json:"name"`
	DefaultHealth float64 `json:"default_health"`
	Health float64 `json:"health"`
	Hunger float64 `json:"hunger"`
	Thirst float64 `json:"thirst"`
	Rubles int `json:"rubles"`
	Inventory []InventoryItem
	Weight float64 `json:"weight"`
	CarryingCapacity float64 `json:"carrying_capacity"`
	Stamina float64 `json:"stamina"`
	Location Coordinate `json:"location"`
	IndoorLocation IndoorCoordinate `json:"indoor_location"`

	LandSpeed float64 `json:"land_speed"`
	SnowSpeed float64 `json:"snow_speed"`
	ClimbSpeed float64 `json:"climb_speed"`

	Aim float64 `json:"aim"`
	Dodge float64 `json:"dodge"`

	Vision float64 `json:"vision"`

	Type CharType `json:"type"`
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
}


type Clothing struct {
	Name string `json:"name"`
	ColdProtection float64 `json:"cold_protection"`
	HeatProtection float64 `json:"heat_protection"`
	Bulletproof float64 `json:"bulletproof"`
	MeleeProtection float64 `json:"melee_protection"`
	Weight float64 `json:"weight"`
	Durability float64 `json:"durability"`
	BodyPart BodyArmor `json:"body_part"`
	Modifier Modifier `json:"modifier"`
}

type Weapon interface {
	attack(attacker *Character, defender *Character)
	estimateHitChance(bodyPart BodyArmor,attacker *Character, defender *Character)float64
}


type Gun struct{
	Name string `json:"name"`
	Calibre string `json:"calibre"`
	EffectiveRange float64 `json:"effective_range"`
	MaxDamage float64 `json:"damage"`
	Accuracy float64 `json:"accuracy"`
	Recoil float64 `json:"recoil"`
	TimeBetweenShots float64 `json:"time_between_shots"`
	ReloadTime float64 `json:"reload_time"`
	LoadedMagazine Ammo `json:"loaded_magazine"`
	Weight float64 `json:"weight"`
	Durability float64 `json:"durability"`
}

type Ammo struct {
	Name string `json:"name"`
	Calibre string `json:"calibre"`
	ArmorPiercing float64 `json:"armor_piercing"`
}

type Melee struct{
	Name string `json:"name"`
	Range float64 `json:"range"`
	PiercingDamage float64 `json:"piercing_damage"`
	SlashingDamage float64 `json:"slashing_damage"`
	CrushingDamage float64 `json:"crushing_damage"`
	Weight float64 `json:"weight"`
	Durability float64 `json:"durability"`
}

type Coordinate struct{
	Row int `json:"row"`
	Col int `json:"col"`
}

type IndoorCoordinate struct {
	MapLocation Coordinate `json:"location"`
	Floor int `json:"floor"`
	RoomLocation Coordinate `json:"room_location"`
}

func (c Coordinate)isNull()bool{
	return c.Col==-1 || c.Row==-1
}

func (p Character)isIndoor()bool{
	return p.Location.isNull()
}

type Place struct {

}

//item price evaluations
func (m Modifier) evaluate()float64{
	if m.IsPercent{
		return 1+(m.Modifier/100)
	}
	//modify based on usefulness of stat
	return 1.5

}

func (g Gun)evaluate()int{
	rangeScore := g.EffectiveRange
	damageScore := g.MaxDamage
	accuracyScore := 2.5*(g.Accuracy/100)
	ROFScore := (MAX_GUN_TIME_BETWEEN_SHOTS+0.2*MAX_GUN_TIME_BETWEEN_SHOTS-g.TimeBetweenShots)/MAX_GUN_TIME_BETWEEN_SHOTS
	recoilScore := ((110-g.Recoil)/100)/ROFScore
	reloadScore := (MAX_GUN_RELOAD_TIME+0.2*MAX_GUN_RELOAD_TIME-g.ReloadTime)/MAX_GUN_RELOAD_TIME
	weightScore := (MAX_GUN_WEIGHT+0.2*MAX_GUN_WEIGHT-g.Weight)/MAX_GUN_WEIGHT
	durabilityScore := g.Durability/100
	price := int(rangeScore*damageScore*accuracyScore*ROFScore*recoilScore*reloadScore*weightScore*durabilityScore)
	if LOG_MODE >=1{
		fmt.Printf("Valued %s at %d%c\n",g.Name,price,RUBLE)
	}
	if LOG_MODE==DEBUG{
		fmt.Printf("Range score: %f\n",rangeScore)
		fmt.Printf("Damage score: %f\n",damageScore)
		fmt.Printf("Accuracy score: %f\n",accuracyScore)
		fmt.Printf("Recoil score: %f\n",recoilScore)
		fmt.Printf("Rate of fire score: %f\n",ROFScore)
		fmt.Printf("Reload speed score: %f\n",reloadScore)
		fmt.Printf("Weight score: %f\n",weightScore)
		fmt.Printf("Durability score: %f\n",durabilityScore)
		fmt.Println()
	}
	return price
}

func (g Gun)estimateHitChance(bodyPart BodyArmor,attacker *Character, defender *Character)float64{
	aimModifier := attacker.Aim/100
	if attacker.isIndoor()!=defender.isIndoor(){
		return 0.0
	}
	distance := calculateDistance(attacker.Location,defender.Location,attacker.isIndoor())
	distanceModifier := getTargetValueNoDir(0,distance,(BULLET_MISS_HARSHNESS*distance)/g.EffectiveRange,false,g.EffectiveRange)
	durabilityModifier := g.Durability/100
	bodyPartModifier := BODY_HIT_MODIFIERS[bodyPart]
	probability := aimModifier*distanceModifier*durabilityModifier*bodyPartModifier
	if LOG_MODE>=1{
		fmt.Printf("Shot success probability: %f\n",probability)
	}
	if LOG_MODE==DEBUG{
		fmt.Printf("Aim modifier: %f\n",aimModifier)
		fmt.Printf("Distance: %f\n",distance)
		fmt.Printf("Distance modifier: %f\n",distanceModifier)
		fmt.Printf("Durability modifier: %f\n",durabilityModifier)
		fmt.Printf("Body part modifier: %f\n",bodyPartModifier)
	}
	if probability>=0.99{
		probability = 0.99
	}
	return probability

}

func (c Clothing)evaluate()int{
	coldScore := 3.0*(c.ColdProtection/100)
	heatScore := c.HeatProtection/100
	bulletproofScore := 2.0*(c.Bulletproof/100)
	meleeProtectionScore := c.MeleeProtection/100
	weightScore :=  (BODY_MAX_WEIGHTS[c.BodyPart]+0.2*BODY_MAX_WEIGHTS[c.BodyPart]-c.Weight)/BODY_MAX_WEIGHTS[c.BodyPart]
	priceMod := BODY_SCORES[c.BodyPart]
	modifierScore := c.Modifier.evaluate()
	durabilityScore := c.Durability/100
	price := int(coldScore*heatScore*bulletproofScore*meleeProtectionScore*weightScore*priceMod*modifierScore*durabilityScore*50000)
	if LOG_MODE >=1{
		fmt.Printf("Valued %s at %d%c\n",c.Name,price,RUBLE)
	}
	if LOG_MODE==DEBUG{
		fmt.Printf("Cold score: %f\n",coldScore)
		fmt.Printf("Heat score: %f\n",heatScore)
		fmt.Printf("Bulletproof score: %f\n",bulletproofScore)
		fmt.Printf("Melee protection score: %f\n",meleeProtectionScore)
		fmt.Printf("Body placement modifier: %f\n",priceMod)
		fmt.Printf("Modifier score: %f\n",modifierScore)
		fmt.Printf("Weight score: %f\n",weightScore)
		fmt.Printf("Durability score: %f\n",durabilityScore)
		fmt.Println()
	}
	return price
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



