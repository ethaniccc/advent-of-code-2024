package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	DirectionDown Direction = iota
	DirectionUp
	DirectionLeft
	DirectionRight
	DirectionDownRight
	DirectionDownLeft
	DirectionUpLeft
	DirectionUpRight
	directionCount
)

var (
	wordSearch = make(map[Vec2]byte)
	endBoundry = Vec2{lineLength}

	lineLength = 140
)

type Direction byte

func (d Direction) Modifier() Vec2 {
	switch d {
	case DirectionDown:
		return Vec2{0, -1}
	case DirectionUp:
		return Vec2{0, 1}
	case DirectionLeft:
		return Vec2{-1, 0}
	case DirectionRight:
		return Vec2{1, 0}
	case DirectionDownRight:
		return Vec2{1, -1}
	case DirectionDownLeft:
		return Vec2{-1, -1}
	case DirectionUpRight:
		return Vec2{1, 1}
	case DirectionUpLeft:
		return Vec2{-1, 1}
	default:
		panic("wt(heck) are you????")
	}
}

type Vec2 [2]int

func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v[0] + v2[0], v[1] + v2[1]}
}

func parseInput() {
	wordSearch = make(map[Vec2]byte)
	dat, err := os.ReadFile("input")
	if err != nil {
		panic(err)
	}

	yAxis := 1
	for _, line := range strings.Split(string(dat), "\n") {
		for xAxis := 1; xAxis <= lineLength; xAxis++ {
			wordSearch[Vec2{xAxis, yAxis}] = line[xAxis-1]
		}
		yAxis++
	}
	endBoundry[1] = yAxis - 1
}

func searchXMAS(start Vec2, d Direction) bool {
	fSeq := []byte("XMAS")
	rSeq := []byte("SAMX")

	if sChar, ok := wordSearch[start]; !ok || (sChar != fSeq[0] && sChar != fSeq[1]) {
		return false
	}

	sequence := make([]byte, 4)
	currentPos := start
	for i := 0; i < 4; i++ {
		char, ok := wordSearch[currentPos]
		if !ok {
			return false
		}

		sequence[i] = char
		currentPos = currentPos.Add(d.Modifier())
	}

	if sequenceMatches(sequence, fSeq, rSeq) {
		return true
	}
	return false
}

func searchXShapedMAS(start Vec2) bool {
	fSeq := []byte("MAS")
	rSeq := []byte("SAM")

	startChar, ok := wordSearch[start]
	if !ok || startChar != fSeq[1] {
		return false
	}

	s1, s2 := []byte{
		wordSearch[start.Add(Vec2{-1, 1})],
		wordSearch[start],
		wordSearch[start.Add(Vec2{1, -1})],
	}, []byte{
		wordSearch[start.Add(Vec2{-1, -1})],
		wordSearch[start],
		wordSearch[start.Add(Vec2{1, 1})],
	}

	if sequenceMatches(s1, fSeq, rSeq) && sequenceMatches(s2, fSeq, rSeq) {
		return true
	}
	return false
}

func sequenceMatches(s, fSeq, rSeq []byte) bool {
	return slices.Equal(s, fSeq) || slices.Equal(s, rSeq)
}

func main() {
	parseInput()

	var xmasScore, masXScore int
	for xAxis := 1; xAxis <= endBoundry[0]; xAxis++ {
		for yAxis := 1; yAxis <= endBoundry[1]; yAxis++ {
			searchPosition := Vec2{xAxis, yAxis}
			for currentDirection := DirectionDown; currentDirection < directionCount; currentDirection++ {
				if searchXMAS(searchPosition, currentDirection) {
					xmasScore++
				}
			}

			if searchXShapedMAS(searchPosition) {
				masXScore++
			}
		}
	}

	fmt.Println("found XMAS:", xmasScore)
	fmt.Println("found X-shaped MAS:", masXScore)
}
