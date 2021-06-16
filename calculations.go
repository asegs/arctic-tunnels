package main

import (
	"fmt"
	"math"
)

const OUTDOOR_COORD_WIDTH = 10
const INDOOR_COORD_WIDTH = 1

func getTargetValueNoDir(expectedValue float64,actualValue float64,harshness float64,test bool,spread float64) float64{
	if LOG_MODE==DEBUG && !test{
		fmt.Printf("Expected value: %f\n",expectedValue)
		fmt.Printf("Actual value: %f\n",actualValue)
		fmt.Printf("Harshness: %f\n",harshness)
		fmt.Printf("Spread: %f\n",spread)
		fmt.Println()
	}
	standardDev := spread/harshness

	if standardDev<1 && standardDev>-1 { standardDev = 1 }

	if harshness==0 { harshness = 0.1 }

	if test {
		result := 1 / (standardDev*math.Sqrt(2*math.Pi))
		if result==0 { return 0.0001 }
		return result
	}

	firstTerm := 1/(standardDev*math.Sqrt(2*math.Pi))
	secondTerm := -math.Pow(actualValue-expectedValue,2)/(2*math.Pow(standardDev,2))
	multiplier := 1/(getTargetValueNoDir(expectedValue,expectedValue,harshness,true,spread))
	return multiplier*firstTerm*math.Exp(secondTerm)
}

func calculateDistance(c1 Coordinate,c2 Coordinate,indoor bool)float64{
	if indoor {
		return math.Sqrt(float64((c1.Col*INDOOR_COORD_WIDTH-c2.Col*INDOOR_COORD_WIDTH)*(c1.Col*INDOOR_COORD_WIDTH-c2.Col*INDOOR_COORD_WIDTH) + (c1.Row*INDOOR_COORD_WIDTH-c2.Row*INDOOR_COORD_WIDTH)*(c1.Row*INDOOR_COORD_WIDTH-c2.Row*INDOOR_COORD_WIDTH)))
	}
	return math.Sqrt(float64((c1.Col*OUTDOOR_COORD_WIDTH-c2.Col*OUTDOOR_COORD_WIDTH)*(c1.Col*OUTDOOR_COORD_WIDTH-c2.Col*OUTDOOR_COORD_WIDTH) + (c1.Row*OUTDOOR_COORD_WIDTH-c2.Row*OUTDOOR_COORD_WIDTH)*(c1.Row*OUTDOOR_COORD_WIDTH-c2.Row*OUTDOOR_COORD_WIDTH)))
}

func intAbs(i int)int{
	if i<0 {
		return i * -1
	}
	return i
}

