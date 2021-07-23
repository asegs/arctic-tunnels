package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const MAJOR_TILES_WIDTH = 100
const MAJOR_TILES_HEIGHT = 20
const TILE_SIZE = 20
const RANDOMNESS = 0.9

type topoWeight struct {
	Height int
	Weight int
}

//.,S,/,M,u,&,,
var SYMBOLS_OUTDOOR = [...]rune{'.','S','/','M','u','&',','}
const SYMBOL_COUNT = len(SYMBOLS_OUTDOOR)
var BASE_FREQUENCIES_OUTDOOR = [SYMBOL_COUNT]float64{1.0,0.4,0.5,0.3,0.2,0.05,0.6}
var LIGHT_SNOW_FREQ_MAP =[]float64{4,9}
var HEAVY_SNOW_FREQ_MAP = []float64{2}
var ICE_FREQ_MAP = []float64{0,7}
var CLIMBABLE_ROCK_FREQ_MAP = []float64{2,6}
var UNCLIMBABLE_ROCK_FREQ_MAP = []float64{4,8}
var LAVA_FREQ_MAP = []float64{8}
var SLUSH_FREQ_MAP = []float64{0}

var allFrequencies = [SYMBOL_COUNT][]float64{LIGHT_SNOW_FREQ_MAP,HEAVY_SNOW_FREQ_MAP,ICE_FREQ_MAP,CLIMBABLE_ROCK_FREQ_MAP,UNCLIMBABLE_ROCK_FREQ_MAP,LAVA_FREQ_MAP,SLUSH_FREQ_MAP}

func printTopoMap(m [][]int){
	for _,row:=range m{
		for _,cell := range row{
			fmt.Print(cell)
		}
		fmt.Println()
	}
}

func cellInbounds(row int,col int,height int,width int)bool{
	return (row<height && row>=0) && (col<width && col>=0)
}

func generateTiledMap(passes int,buildContrastMap bool){
	neighborMoves := []int{-1,0,0,1,1,0,0,-1}
	filename := "Config/topo_map.txt"
	if buildContrastMap{
		filename = "Config/contrast_map.txt"
	}
	configTopoString := ReadToString(filename)
	topoRows := strings.Split(configTopoString,"\n")
	rowCount := len(topoRows)
	if rowCount == 0{
		return
	}
	rowLength := len(topoRows[0])
	if rowLength == 0{
		return
	}
	topoMap := make([][]int,rowCount)
	for i:=0;i<rowCount;i++{
		topoMap[i] = make([]int,rowLength)
	}
	for i,row := range topoRows {
		for b,r := range row{
			topoMap[i][b] = ctoiSafe(rune(r))
		}
	}

	contrastMap := make([][]int,rowCount)
	for i:=0;i<rowCount;i++{
		contrastMap[i] = make([]int,rowLength)
		if buildContrastMap{
			for b:=0;b<rowLength;b++{
				contrastMap[i][b] = 1
			}
		}
	}

	if !buildContrastMap {
		contrastConfigString := ReadToString("Config/smooth_contrast_map.txt")
		contrastRows := strings.Split(contrastConfigString,"\n")
		for i,row:=range contrastRows{
			for b,char:=range row{
				contrastMap[i][b] = ctoiSafe(char)
			}
		}
	}

	for passes>0{
		for i,row := range topoMap{
			for b,num := range row {
				for j := 0;j<8;j+=2{
					neighborRow := i+neighborMoves[j]
					neighborCol := b+neighborMoves[j+1]
					if cellInbounds(neighborRow,neighborCol,rowCount,rowLength){
						neighborNum := topoMap[neighborRow][neighborCol]
						contrast := contrastMap[i][b]
						if contrastMap[i][b]<intAbs(num-neighborNum){
							if neighborNum <= num {
								topoMap[neighborRow][neighborCol] = num-contrast
							}else{
								topoMap[i][b] = neighborNum-contrast
							}
						}
					}
				}
			}
		}
		passes--
	}
	if buildContrastMap {
		sb := strings.Builder{}
		for _,row:=range topoMap{
			for _,col:=range row{
				sb.WriteRune(rune(col+48))
			}
			sb.WriteRune('\n')
		}

		Write("Config/smooth_contrast_map.txt",sb.String())
		return
	}
	hugeTopoMap := make([][]rune,rowCount*TILE_SIZE)
	for i:=0;i<rowCount*TILE_SIZE;i++{
		hugeTopoMap[i] = make([]rune,rowLength*TILE_SIZE)
	}
	detailTopoMap := make([][][][]float64,rowCount)
	detailTerrainMap := make([][][][]rune,rowCount)
	for i:=0;i<rowCount;i++{
		detailTopoMap[i] = make([][][]float64,rowLength)
		detailTerrainMap[i] = make([][][]rune,rowLength)
		for b:=0;b<rowLength;b++{
			detailTopoMap[i][b] = make([][]float64,TILE_SIZE)
			detailTerrainMap[i][b] = make([][]rune,TILE_SIZE)
			for n:=0;n<TILE_SIZE;n++{
				detailTopoMap[i][b][n] = make([]float64,TILE_SIZE)
				detailTerrainMap[i][b][n] = make([]rune,TILE_SIZE)
				for z:=0;z<TILE_SIZE;z++{
					detailTopoMap[i][b][n][z] = 0.0
					detailTerrainMap[i][b][n][z] = ' '
				}
			}
		}
	}


	//prepare for DEEP nesting
	//in small tile, each individual cell elevation is: weighted (by distance) average elevation of nearby major tiles, top,right,bottom,left,center (actual center tile of center)
	for i,row := range topoMap{
		for b,cell := range row {
			used := 0.0
			cells := make([]topoWeight,4)
			for j := 0;j<8;j+=2 {
				neighborRow := i + neighborMoves[j]
				neighborCol := b + neighborMoves[j+1]
				if cellInbounds(neighborRow, neighborCol, rowCount, rowLength) {
					cells[(j/2)].Height = topoMap[i][b]
					cells[(j/2)].Weight = 1
					used++
				}
			}
			//math.abs of (n-halfway size * weight) (to check for off map) to get which is closer, then add that number times the height of respective terrain
			//add (1- math.abs(n-10))/size * self

			//terrain too regular, uncommon stuff less likely.  introduce some serious randomness not to score but to if something can spawn at all, for example base value is what it needs to be lower than to spawn
			for n:=0;n<TILE_SIZE;n++{
				for z:=0;z<TILE_SIZE;z++{
					halfwayDistance := 0.5 * TILE_SIZE
					eachTileWeight := 1.0/3.0
					horizontal := math.Max(float64(z)/float64(TILE_SIZE) * float64(cells[1].Height) * float64(cells[1].Weight),(1-float64(z)/float64(TILE_SIZE)) * float64(cells[3].Height) * float64(cells[3].Weight)) * eachTileWeight
					verticalHeight := math.Max(1-((float64(n)-halfwayDistance)/float64(TILE_SIZE)) * float64(cells[0].Height) * float64(cells[0].Weight),(float64(n)-halfwayDistance)/float64(TILE_SIZE) * float64(cells[2].Height) * float64(cells[2].Weight)) * eachTileWeight
					horizontalCenter := 1 - math.Abs(halfwayDistance-float64(n)) * 0.5
					verticalCenter := 1- math.Abs(halfwayDistance-float64(z)) * 0.5
					centerHeight := float64(cell) * eachTileWeight * (horizontalCenter + verticalCenter)

					indivHeight := horizontal + verticalHeight + centerHeight

					detailTopoMap[i][b][n][z] = indivHeight
					scores := make([]float64,SYMBOL_COUNT)
					for x:=0;x<SYMBOL_COUNT;x++{
						places := allFrequencies[x]
						baseFreq := BASE_FREQUENCIES_OUTDOOR[x]
						toScore := make([]float64,len(places))
						for m:=0;m<len(places);m++{
							toScore[m] = getTargetValueNoDir(0,math.Abs(places[m]-indivHeight),5.0,false,3.0)
						}
						scores[x],_ = floatMax(toScore)
						randomMultiplier := 1.0 - RANDOMNESS
						if r1.Float64()<baseFreq{
							randomMultiplier = 1.0
						}
						scores[x] = scores[x] * randomMultiplier
					}
					_,index := floatMax(scores)
					detailTerrainMap[i][b][n][z] = SYMBOLS_OUTDOOR[index]
					hugeTopoMap[i*TILE_SIZE+n][b*TILE_SIZE+z] = SYMBOLS_OUTDOOR[index]
				}
			}
		}
	}
	sb := strings.Builder{}
	for _,row := range hugeTopoMap {
		for _,cell := range row {
			sb.WriteRune(cell)
		}
		sb.WriteRune('\n')
	}
	Write("Config/total_map.txt",sb.String())
}

func main(){
	start := time.Now()
	generateTiledMap(10,false)
	end := time.Now()
	fmt.Println(end.Sub(start))
}