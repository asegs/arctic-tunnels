package main

import (
	"fmt"
	"math"
	"math/rand"
)

const MAX_GUN_WEIGHT = 24
const MAX_GUN_RELOAD_TIME = 4000
const MAX_GUN_TIME_BETWEEN_SHOTS = 1500
const BULLET_MISS_HARSHNESS = 1.75

var BODY_SCORES = [...]float64{1.5,1.0,0.25,0.35,0.45}
var BODY_MAX_WEIGHTS = [...]float64{8,40,1.5,6,5}
var BODY_HIT_MODIFIERS = [...]float64{0.5,2.0,0.65,1.25,1.0}
var BODY_DAMAGE_MODIFIERS = [...]float64{3.0,1.25,0.6,0.5,0.35}

var GUN_NAME_MODIFIERS = [...]string{"Short-Barrel","Long-Barrel","Light","Hard-Hitting","Imprecise","Sniper","Jerky","Smooth","Easy-Use","Stiff","Light","Sturdy"}
var CONDITIONS = [...]string{"Mint","Well Cared For","Solid","Rusty","Beater","Broken"}

//calibre information
var CALIBRES = [...]string{".22 LR","9mm","5.56 NATO","7.62x39","7.62 NATO",".45 ACP",".338 Lapua",".50 BMG","6.5 Creedmoor"}
const l = len(CALIBRES)
var CALIBRE_BASE_RANGE = [l]float64{75.0,100.0,220.0,280.0,500.0,120.0,1200.0,1750.0,500.0}
var CALIBRE_BASE_DAMAGE = [l]float64{15.0,45.0,68.0,77.0,105.0,60.0,155.0,210.0,95.0}
var CALIBRE_RECOIL_LEVEL = [l]float64{3.0,25.0,31.0,55.0,72.0,52.0,90.0,100.0,39.0}

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

func (g Gun)calculateDamage(bodyPart BodyArmor,attacker *Character, defender *Character)float64{
	dmgModifier := BODY_DAMAGE_MODIFIERS[bodyPart]
	distance := calculateDistance(attacker.Location,defender.Location,attacker.isIndoor())
	distanceModifier := getTargetValueNoDir(0,distance,(BULLET_MISS_HARSHNESS*distance)/g.EffectiveRange,false,g.EffectiveRange)
	baseDamage := g.MaxDamage
	armorDurability := defender.Armor[bodyPart].Durability/100
	bulletproofModifier := (100-defender.Armor[bodyPart].Bulletproof*(armorDurability/100))/100
	bulletAppropriateModifier := math.Abs(g.LoadedMagazine.ArmorPiercing-bulletproofModifier*100)/100
	damage := baseDamage*dmgModifier*distanceModifier*bulletproofModifier*bulletAppropriateModifier
	if LOG_MODE>=1{
		fmt.Printf("%s did %f damage to %s",attacker.Name,damage,defender.Name)
	}
	if LOG_MODE==DEBUG{
		fmt.Printf("Body part damage modifier: %f\n",dmgModifier)
		fmt.Printf("Distance: %f\n",distance)
		fmt.Printf("Distance modifier: %f\n",distanceModifier)
		fmt.Printf("Base damage: %f\n",baseDamage)
		fmt.Printf("Armor durability: %f\n",armorDurability)
		fmt.Printf("Bulletproof modifier: %f\n",bulletproofModifier)
		fmt.Printf("Bullet appropriate modifier: %f\n",bulletAppropriateModifier)
	}

	return damage
}

//redo attack methods as attacker.attack(bp,weapon,defender)
func (g Gun)attack(bodyPart BodyArmor,attacker *Character,defender *Character){
	if g.LoadedMagazine.Rounds==0{
		//need to reload!
		return
	}
	hitChance := g.estimateHitChance(bodyPart,attacker,defender)
	roll := rand.Float64()
	if roll<=hitChance{
		damage := g.calculateDamage(bodyPart,attacker,defender)
		g.LoadedMagazine.Rounds--
		defender.Health-=damage
	}
}

//random stats, generate full name based on RNG of these stats
func generateGun(baseName string,calibreIdx int,distanceMod float64,rof int,damageMod float64,condition int,accuracy int)Gun{

}
