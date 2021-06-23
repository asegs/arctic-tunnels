package main

import (
	"fmt"
	"math/rand"
	"time"
)

type LogMode int

type BodyArmor int
type ItemType int
type CharType int

var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)


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


const LOG_MODE = BASIC

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
	Moving bool `json:"moving"`

	Aim float64 `json:"aim"`
	Dodge float64 `json:"dodge"`

	Vision float64 `json:"vision"`

	Type CharType `json:"type"`
	
	Armor [5]Clothing `json:"armor"`
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
	attack(bodyPart BodyArmor, attacker *Character, defender *Character)
	estimateHitChance(bodyPart BodyArmor,attacker *Character, defender *Character)float64
	calculateDamage(bodyPart BodyArmor,attacker *Character,defender *Character)float64
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
	LoadedMagazine Magazine `json:"loaded_magazine"`
	Weight float64 `json:"weight"`
	Durability float64 `json:"durability"`
}

type Magazine struct {
	Rounds int `json:"rounds"`
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

//if shot misses, chance of hitting other part (go down the line, most likely to least likely, each time multiply by 0.5 prob



