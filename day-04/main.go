package main

import (
	"fmt"
	"os"
	"regexp"
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
	wordSearch     = make(map[Vec2]byte)
	foundPositions = make(map[Vec2]struct{})
	searchResults  = make(map[SearchResult]struct{})
	endBoundry     = Vec2{0, 1}

	forwardPattern  = regexp.MustCompile(`^X+MAS$`)
	reversePattern  = regexp.MustCompile(`^S+AMX$`)
	forwardSequence = []byte("XMAS")
	reverseSequence = []byte("SAMX")
	wordLength      = 4
	lineLength      = 140
	actualResults   int
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

func (d Direction) Opposite() Direction {
	switch d {
	case DirectionDown:
		return DirectionUp
	case DirectionUp:
		return DirectionDown
	case DirectionLeft:
		return DirectionRight
	case DirectionRight:
		return DirectionLeft
	case DirectionDownRight:
		return DirectionUpLeft
	case DirectionDownLeft:
		return DirectionUpRight
	case DirectionUpRight:
		return DirectionDownLeft
	case DirectionUpLeft:
		return DirectionDownRight
	default:
		panic("idk wt(heck) is the opposite of you???")
	}
}

type Vec2 [2]int

func (v Vec2) WithinBoundary() bool {
	inBoundsX := v[0] >= 1 && v[0] <= endBoundry[0]
	inBoundsY := v[1] >= 1 && v[1] <= endBoundry[1]

	return inBoundsX && inBoundsY
}

type SearchResult struct {
	Position  Vec2
	Direction Direction
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

	endBoundry[0] = lineLength
	endBoundry[1] = yAxis - 1
}

func trySearchingDirection(start Vec2, d Direction) {
	// First get the sequence of the letters in the direction.
	sequence := make([]byte, wordLength)
	results := make([]SearchResult, wordLength)
	currentPos := start
	modifier := d.Modifier()
	for i := 0; i < wordLength; i++ {
		if !currentPos.WithinBoundary() {
			// Do not attempt this direction if it eventually goes out of bounds of the word search.
			return
		}

		sequence[i] = wordSearch[currentPos]
		results[i] = SearchResult{Position: currentPos, Direction: d}

		currentPos[0] += modifier[0]
		currentPos[1] += modifier[1]
	}

	// Now, check if we already have this search result found from another search.
	if _, ok := searchResults[results[0]]; ok {
		return
	}
	// We also have to make sure none in the opposite direction exist.
	if _, ok := searchResults[SearchResult{Position: start, Direction: d.Opposite()}]; ok {
		return
	}

	// Finally, we check if the result is valid, and if so, add it to the found cache
	if slices.Equal(reverseSequence, sequence) || slices.Equal(forwardSequence, sequence) {
		for _, result := range results {
			searchResults[result] = struct{}{}
			foundPositions[result.Position] = struct{}{}
		}
		actualResults++
		fmt.Printf("%s FOUND at (%v, direction=%d)\n", string(sequence), start, d)
	} else {
		// Check for "overlap" by getting next letter and seeing if any one could be removed to get "XMAS"
		currentPos[0] += modifier[0]
		currentPos[1] += modifier[1]
		sequence = append(sequence, wordSearch[currentPos])
		results = append(results, SearchResult{Position: currentPos, Direction: d})

		if forwardPattern.Match(sequence) || reversePattern.Match(sequence) {
			for _, result := range results {
				searchResults[result] = struct{}{}
				foundPositions[result.Position] = struct{}{}
			}
			actualResults++
			fmt.Printf("%s weird overlap FOUND at (%v, direction=%d)\n", string(sequence), start, d)
		}
	}

}

func main() {
	parseInput()
	for xAxis := 1; xAxis <= endBoundry[0]; xAxis++ {
		for yAxis := 1; yAxis <= endBoundry[1]; yAxis++ {
			searchPosition := Vec2{xAxis, yAxis}
			for currentDirection := DirectionDown; currentDirection < directionCount; currentDirection++ {
				trySearchingDirection(searchPosition, currentDirection)
			}
		}
	}

	/* for yAxis := 1; yAxis <= endBoundry[1]; yAxis++ {
		dat := ""
		for xAxis := 1; xAxis <= endBoundry[0]; xAxis++ {
			pos := Vec2{xAxis, yAxis}
			if _, ok := foundPositions[pos]; ok {
				dat += "*"
			} else {
				dat += string(wordSearch[pos])
			}
		}
		fmt.Println(dat)
	} */

	fmt.Println("found:", actualResults)
}
