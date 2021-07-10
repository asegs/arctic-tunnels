package main

import (
	"fmt"
	"math/rand"
)

type direction struct {
	vert int
	hor int
}

func isPerpendicular(d1 direction, d2 direction)bool{
	return (d1.hor==d2.vert && d1.vert == -d2.hor) || (d1.hor==-d2.vert && d1.vert==d2.hor)
}

func isExiting(row int,col int,selectedDirection direction,height int,width int)bool{
	return (row==0 && selectedDirection.vert==-1) ||
		(col==0 && selectedDirection.hor==-1) ||
		(row==height-1 && selectedDirection.vert==1) ||
		(col==width-1 && selectedDirection.hor==1)
}

func populateDirections()[4]direction{
	var directions [4]direction
	directions[0] = direction{
		vert: -1,
		hor:  0,
	}
	directions[1] = direction{
		vert: 1,
		hor:  0,
	}
	directions[2] = direction{
		vert: 0,
		hor:  -1,
	}
	directions[3] = direction{
		vert: 0,
		hor:  1,
	}
	return directions
}

func createEmptyRoom(height int,width int,maxTunnels int,maxLength int)[][]rune{
	room := make([][]rune,height)
	for i:=0;i<height;i++{
		room[i] = make([]rune,width)
		for b:=0;b<width;b++{
			room[i][b] = '|'
		}
	}

	totalMax := maxTunnels

	for z:=0;z<5;z++{
		maxTunnels = totalMax

		currentRow := rand.Intn(height)
		currentCol := rand.Intn(width)

		fmt.Println(currentCol)
		fmt.Println(currentRow)

		directions := populateDirections()
		lastDirection := direction{
			vert: 1,
			hor:  0,
		}
		var randomDirection direction

		for maxTunnels > 0{
			randomDirection = directions[rand.Intn(4)]
			for !isPerpendicular(randomDirection,lastDirection) {
				randomDirection = directions[rand.Intn(4)]
			}
			randomLength := rand.Intn(maxLength)+1
			tunnelLength := 0
			for tunnelLength<randomLength{
				if isExiting(currentRow,currentCol,randomDirection,height,width){
					break
				}else{
					room[currentRow][currentCol] = ' '
					currentRow += randomDirection.vert
					currentCol += randomDirection.hor
					tunnelLength++
				}
			}
			if tunnelLength>0{
				lastDirection = randomDirection
				maxTunnels--
			}
		}
	}

	return room
}

func displayRoom(room [][]rune){
	for _,row := range room{
		for _,item := range row{
			fmt.Printf("%c",item)
		}
		fmt.Println()
	}
}

//func main(){
//	room := createEmptyRoom(15,100,50,4)
//	displayRoom(room)
//}