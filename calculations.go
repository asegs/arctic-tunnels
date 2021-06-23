package main

import (
	"fmt"
	"math"
)

const OUTDOOR_COORD_WIDTH = 10
const INDOOR_COORD_WIDTH = 1

func getTargetValueNoDir(expectedValue float64,actualValue float64,harshness float64,getMax bool,spread float64) float64{
	if LOG_MODE==DEBUG && !getMax{
		fmt.Printf("Expected value: %f\n",expectedValue)
		fmt.Printf("Actual value: %f\n",actualValue)
		fmt.Printf("Harshness: %f\n",harshness)
		fmt.Printf("Spread: %f\n",spread)
		fmt.Println()
	}
	standardDev := spread/harshness

	if standardDev<1 && standardDev>-1 { standardDev = 1 }

	if harshness==0 { harshness = 0.1 }

	if getMax {
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

func deviationFromCenter(lowerBound float64,upperBound float64,pick float64)float64{
	return (pick-(lowerBound+upperBound)/2)/(upperBound-lowerBound)
}

func pickRandomInRange(lowerBound float64,upperBound float64)float64{
	return r1.Float64()*(upperBound-lowerBound)+lowerBound
}

func pickRandomVariedAround(center float64,variability float64)float64{
	return pickRandomInRange(getAroundRange(center,variability))
}

func getAroundRange(center float64,variability float64)(float64,float64){
	return center*(1-variability),center*(1+variability)
}

func MaxOf(vars ...float64) float64 {
	max := math.Abs(vars[0])

	for _, i := range vars {
		if max < math.Abs(i) {
			max = math.Abs(i)
		}
	}

	return max
}

func signToOne(num float64)int{
	return int(num/math.Abs(num))
}

func MaxIndex(m float64,values ... float64)int{
	for i,value := range values {
		absValue := math.Abs(value)
		if absValue==m{
			if absValue != value {
				return i*2
			}
			return i*2+1
		}
	}
	return -1
}
