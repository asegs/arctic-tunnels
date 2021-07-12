package main

import (
	"fmt"
	"strings"
)

const MAJOR_TILES_WIDTH = 100
const MAJOR_TILES_HEIGHT = 20
const TILE_SIZE = 20

//.,S,/,M,u,&,,
var BASE_FREQUENCIES_OUTDOOR = [...]float64{1.0,0.4,0.5,0.3,0.2,0.05,0.6}
var LIGHT_SNOW_FREQ_MAP =[]float64{4,9}
var HEAVY_SNOW_FREQ_MAP = []float64{2}
var ICE_FREQ_MAP = []float64{0,7}
var CLIMBABLE_ROCK_FREQ_MAP = []float64{2,6}
var UNCLIMBABLE_ROCK_FREQ_MAP = []float64{4,8}
var LAVA_FREQ_MAP = []float64{8}
var SLUSH_FREQ_MAP = []float64{0}

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
		fmt.Printf("Pass completed, %d passes remaining\n",passes)
	}
	printTopoMap(topoMap)
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
	detailTopoMap := make([][][][]float64,rowCount)
	for i:=0;i<rowCount;i++{
		detailTopoMap[i] = make([][][]float64,rowLength)
		for b:=0;b<rowLength;b++{
			detailTopoMap[i][b] = make([][]float64,TILE_SIZE)
			for n:=0;n<TILE_SIZE;n++{
				detailTopoMap[i][b][n] = make([]float64,TILE_SIZE)
				for z:=0;z<TILE_SIZE;z++{
					detailTopoMap[i][b][n][z] = 0.0
				}
			}
		}
	}
	//prepare for DEEP nesting
	//in small tile, each individual cell elevation is: weighted (by distance) average elevation of nearby major tiles, top,right,bottom,left,center (actual center tile of center)
	for i,row := range topoMap{
		for b,cell := range row {
			for j := 0;j<8;j+=2 {
				neighborRow := i + neighborMoves[j]
				neighborCol := b + neighborMoves[j+1]
				if cellInbounds(neighborRow, neighborCol, rowCount, rowLength) {

				}
			}
		}
	}
}

func main(){
	generateTiledMap(10,false)
}