package main

import (
	"fmt"
	"time"
)

func main(){
	ar15 := Gun{
		Name:             "AR15",
		Calibre:          "5.56 NATO",
		EffectiveRange:   120,
		MaxDamage:        45,
		Accuracy:         80,
		Recoil:           15,
		TimeBetweenShots: 250,
		ReloadTime:       2000,
		LoadedMagazine:   Magazine{},
		Weight:           12,
		Durability:       85,
	}
	ar15.evaluate()

	ak47 := Gun{
		Name:             "AK47",
		Calibre:          "7.62x39",
		EffectiveRange:   100,
		MaxDamage:        100,
		Accuracy:         50,
		Recoil:           80,
		TimeBetweenShots: 100,
		ReloadTime:       1600,
		LoadedMagazine:   Magazine{},
		Weight:           10,
		Durability:       85,
	}

	AWP := Gun{
		Name:             "AWP",
		Calibre:          ".338 Lapua",
		EffectiveRange:   1500,
		MaxDamage:        200,
		Accuracy:         100,
		Recoil:           100,
		TimeBetweenShots: 1200,
		ReloadTime:       3000,
		LoadedMagazine:   Magazine{},
		Weight:           20,
		Durability:       100,
	}

	ak47.evaluate()
	AWP.evaluate()

	tshirt := Clothing{
		Name:            "T-Shirt",
		ColdProtection:  8,
		HeatProtection:  3,
		Bulletproof:     1,
		MeleeProtection: 12,
		Weight:          0.5,
		Durability:      85,
		BodyPart:        CHEST,
		Modifier:        Modifier{
			StatEffected: "AIM",
			Modifier:     50,
			IsPercent:    true,
		},
	}
	tshirt.evaluate()

	exoskeleton := Clothing{
		Name:            "Exoskeleton Suit",
		ColdProtection:  45,
		HeatProtection:  80,
		Bulletproof:     100,
		MeleeProtection: 100,
		Weight:          40,
		Durability:      90,
		BodyPart:        CHEST,
		Modifier:        Modifier{},
	}
	exoskeleton.evaluate()

	a := [5]Clothing{{},exoskeleton,{},{},{}}
	t := [5]Clothing{{},tshirt,{},{},{}}

	attacker := Character{
		Name:             "Attacker",
		DefaultHealth:    0,
		Health:           0,
		Hunger:           0,
		Thirst:           0,
		Rubles:           0,
		Inventory:        nil,
		Weight:           0,
		CarryingCapacity: 0,
		Stamina:          0,
		Location:         Coordinate{
			Row: 5,
			Col: 5,
		},
		IndoorLocation:   IndoorCoordinate{},
		LandSpeed:        0,
		SnowSpeed:        0,
		ClimbSpeed:       0,
		Aim:              80,
		Dodge:            0,
		Vision:           0,
		Type:             0,
		Armor: t,
	}

	defender := Character{
		Name:             "Defender",
		DefaultHealth:    0,
		Health:           0,
		Hunger:           0,
		Thirst:           0,
		Rubles:           0,
		Inventory:        nil,
		Weight:           0,
		CarryingCapacity: 0,
		Stamina:          0,
		Location:         Coordinate{
			Row: 12,
			Col: 16,
		},
		IndoorLocation:   IndoorCoordinate{},
		LandSpeed:        0,
		SnowSpeed:        0,
		ClimbSpeed:       0,
		Aim:              0,
		Dodge:            0,
		Vision:           0,
		Type:             0,
		Armor: a,
	}

	ak := gunCreateStandardAK()
	ak.LoadedMagazine.Rounds+=1000
	fmt.Println(ak)
	s := time.Now()
	for i:=0;i<1000;i++{
		ak.attack(CHEST,&attacker,&defender)
	}
	e := time.Now()
	fmt.Println(e.Sub(s))
	fmt.Println(defender.Health)
}
