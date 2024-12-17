package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Vec2 [2]int

func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v[0] + v2[0], v[1] + v2[1]}
}

var (
	tMap           = make(map[Vec2]int64)
	startPositions = []Vec2{}

	scores = make(map[Vec2]int)
)

func parseInput() {
	dat, err := os.ReadFile("input")
	if err != nil {
		panic(err)
	}

	newline := "\n"
	if runtime.GOOS == "windows" {
		newline = "\r"
	}

	lines := strings.Split(string(dat), newline)
	for yAxis, line := range lines {
		for xAxis := 0; xAxis < len(line); xAxis++ {
			num, err := strconv.ParseInt(string(line[xAxis]), 10, 64)
			if err != nil {
				panic(err)
			}

			pos := Vec2{xAxis + 1, yAxis + 1}
			tMap[pos] = num
			if num == 0 {
				startPositions = append(startPositions, pos)
			}
		}
	}
}

func test(pos Vec2, expected int64) bool {
	if n, ok := tMap[pos]; ok && n == expected {
		return true
	}
	return false
}

func branch(pos Vec2, current int64, found map[Vec2]struct{}) {
	if current == 10 {
		found[pos] = struct{}{}
		if _, ok := scores[pos]; !ok {
			scores[pos] = 1
		} else {
			scores[pos]++
		}
		return
	}

	up := pos.Add(Vec2{0, -1})
	down := pos.Add(Vec2{0, 1})
	left := pos.Add(Vec2{-1, 0})
	right := pos.Add(Vec2{1, 0})

	if test(up, current) {
		branch(up, current+1, found)
	}
	if test(down, current) {
		branch(down, current+1, found)
	}
	if test(left, current) {
		branch(left, current+1, found)
	}
	if test(right, current) {
		branch(right, current+1, found)
	}
}

func main() {
	parseInput()

	score := 0
	for _, origin := range startPositions {
		found := make(map[Vec2]struct{})
		branch(origin, 1, found)
		score += len(found)
	}
	fmt.Println(score)

	totalScore := 0
	for _, s := range scores {
		totalScore += s
	}
	fmt.Println(totalScore)
}
