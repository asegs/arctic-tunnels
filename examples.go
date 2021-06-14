package main

func show(){
	ar15 := Gun{
		Name:             "AR15",
		Calibre:          "5.56 NATO",
		EffectiveRange:   200,
		MaxDamage:        45,
		Accuracy:         80,
		Recoil:           15,
		TimeBetweenShots: 250,
		ReloadTime:       2000,
		LoadedMagazine:   Ammo{},
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
		LoadedMagazine:   Ammo{},
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
		LoadedMagazine:   Ammo{},
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
		Durability:      100,
		BodyPart:        CHEST,
		Modifier:        Modifier{},
	}
	exoskeleton.evaluate()
}
