package main

import (
	"fmt"
)

func main(){

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
		Armor: a,
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
		Armor: t,
	}
	fmt.Println(attacker)
	fmt.Println(defender)
	ak := gunCreateStandardAK()
	fmt.Println(ak)
	ak.evaluate()
	ar := gunCreateStandardAR()
	fmt.Println(ar)
	ar.evaluate()
	sni := gunCreateStandardSniper()
	fmt.Println(sni)
	sni.evaluate()
	sho := gunCreateStandardShotgun()
	fmt.Println(sho)
	sho.evaluate()
}
