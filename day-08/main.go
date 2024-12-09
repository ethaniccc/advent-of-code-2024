package main

import (
	"fmt"
	"os"
	"strings"
)

type Vec2 [2]int

func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v[0] + v2[0], v[1] + v2[1]}
}

func (v Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v[0] - v2[0], v[1] - v2[1]}
}

func (v Vec2) InBounds() bool {
	return v[0] >= 1 && v[1] >= 1 && v[0] <= maxBoundry[0] && v[1] <= maxBoundry[1]
}

var (
	antinodes  = make(map[Vec2]struct{})
	antennas   = make(map[Vec2]byte)
	maxBoundry Vec2
)

func parseInput() {
	dat, err := os.ReadFile("input")
	if err != nil {
		panic(err)
	}

	expectedLineLength := 0
	lineCount := 0
	nonAntena := byte('.')
	for _, line := range strings.Split(string(dat), "\n") {
		if lineLength := len(line); expectedLineLength != 0 && lineLength != expectedLineLength {
			continue
		} else if expectedLineLength == 0 {
			expectedLineLength = lineLength
			maxBoundry[0] = expectedLineLength
		}
		lineCount++
		maxBoundry[1] = lineCount

		for index := 0; index < expectedLineLength; index++ {
			char := line[index]
			if char != nonAntena {
				antennas[Vec2{index + 1, lineCount}] = char
			}
		}
	}
}

func main() {
	parseInput()

	for antennaPos, currentAntenna := range antennas {
		for otherAntennaPos, otherAntenna := range antennas {
			if otherAntennaPos == antennaPos || currentAntenna != otherAntenna {
				continue
			}

			delta := antennaPos.Sub(otherAntennaPos)
			for pos1 := otherAntennaPos; pos1.InBounds(); pos1 = pos1.Sub(delta) {
				antinodes[pos1] = struct{}{}
			}
			for pos2 := antennaPos; pos2.InBounds(); pos2 = pos2.Add(delta) {
				antinodes[pos2] = struct{}{}
			}
		}
	}

	fmt.Println(len(antinodes))
}
