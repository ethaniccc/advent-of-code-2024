package main

import (
	"fmt"
	"maps"
	"os"
	"strings"
	"sync"
)

type Direction byte

func (d Direction) Step() Vec2 {
	switch d {
	case DirectionUp:
		return Vec2{0, -1}
	case DirectionDown:
		return Vec2{0, 1}
	case DirectionLeft:
		return Vec2{-1, 0}
	case DirectionRight:
		return Vec2{1, 0}
	default:
		panic("this should never happen...")
	}
}

func (d Direction) RotateRight() Direction {
	if d == DirectionLeft {
		return DirectionUp
	}
	return d + 1
}

type Vec2 [2]int

func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v[0] + v2[0], v[1] + v2[1]}
}

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
	directionCount
)

var (
	maxBoundry Vec2
	guardStart Vec2

	initalObstacles = make(map[Vec2]struct{})
	loopObstacles   = make(map[Vec2]struct{})
	visited         = make(map[Vec2]struct{})

	loMu, vMu sync.RWMutex
)

func parseInput() {
	dat, err := os.ReadFile("input")
	if err != nil {
		panic(err)
	}

	expectedLineLength, yAxis := 0, 0
	obstacle := byte('#')
	guard := byte('^')

	for _, line := range strings.Split(string(dat), "\n") {
		lineLength := len(line)
		if expectedLineLength == 0 {
			expectedLineLength = lineLength
		} else if expectedLineLength != lineLength {
			continue
		}

		for xAxis := 0; xAxis < lineLength; xAxis++ {
			if char := line[xAxis]; char == obstacle {
				initalObstacles[Vec2{xAxis, yAxis}] = struct{}{}
			} else if char == guard {
				guardStart = Vec2{xAxis, yAxis}
				fmt.Println("found guard at", guardStart)
			}
		}
		yAxis++
	}

	if guardStart == (Vec2{}) {
		panic("no guard found")
	}
	maxBoundry = Vec2{expectedLineLength, yAxis - 1}
}

func traverseUntilLoopOrEnd(tempObstacle *Vec2, obstacles map[Vec2]struct{}) (looped bool) {
	if tempObstacle != nil {
		// Check to see if we've already determined the temporary obstacle to be a cause for the guard to loop.
		loMu.RLock()
		_, tempObFound := loopObstacles[*tempObstacle]
		loMu.RUnlock()
		if tempObFound {
			return false
		}

		// Only add a temporary obstacle if there isn't an existing obstacle from the original input.
		if _, ok := obstacles[*tempObstacle]; !ok {
			obstacles[*tempObstacle] = struct{}{}
			defer delete(obstacles, *tempObstacle)
		} else {
			// If tempObstacle is not nil, we are checking for a potential loop. If the obstacle already exists in
			// the obstacle map, it is impossible for it to loop.
			return false
		}
	}

	collisions := make(map[Direction]map[Vec2]struct{})
	dueToTempObstacle := false

	for d := DirectionUp; d < directionCount; d++ {
		collisions[d] = make(map[Vec2]struct{})
	}

	guardDirection := DirectionUp
	guardPos := guardStart

	// Keep taking another step as long as there is no obstacle in the way.
	for {
		step := guardDirection.Step()
		nextPos := guardPos.Add(step)
		if tempObstacle == nil {
			visited[guardPos] = struct{}{}
		}

		for {
			if outOfBounds(nextPos) {
				// If the position is no longer in bounds, we can stop searching for possible positions.
				return false
			} else if _, ok := obstacles[nextPos]; ok {
				// If this is to test to see if an obstacle can cause the guard to loop, check if
				// the guard is going in a loop.
				if tempObstacle != nil {
					if !dueToTempObstacle && nextPos == *tempObstacle {
						dueToTempObstacle = true
					}
					if _, ok := collisions[guardDirection][nextPos]; ok {
						return dueToTempObstacle
					}
					collisions[guardDirection][nextPos] = struct{}{}
				}

				// This is reached when the guard encounters an obstacle but is still in bounds. We rotate the direction to the right.
				guardDirection = guardDirection.RotateRight()
				break
			}

			guardPos = nextPos
			nextPos = nextPos.Add(step)
			if tempObstacle == nil {
				visited[guardPos] = struct{}{}
			}
		}
	}
}

func tryPossibleLoops(pos Vec2) {
	obstacles := maps.Clone(initalObstacles)
	for _, obPos := range []Vec2{
		pos.Add(DirectionUp.Step()),
		pos.Add(DirectionDown.Step()),
		pos.Add(DirectionLeft.Step()),
		pos.Add(DirectionRight.Step()),
	} {
		// Only add the position to the temp obstacles map if there isn't already an existing obstacle.
		if looped := traverseUntilLoopOrEnd(&obPos, obstacles); looped {
			loMu.Lock()
			loopObstacles[obPos] = struct{}{}
			loMu.Unlock()
		}
	}
}

func outOfBounds(pos Vec2) bool {
	return (pos[0] < 0 || pos[0] > maxBoundry[0]) || (pos[1] < 0 || pos[1] > maxBoundry[1])
}

func main() {
	parseInput()
	traverseUntilLoopOrEnd(nil, initalObstacles)
	fmt.Println("guard traveled to", len(visited), "unique positions")

	var wg sync.WaitGroup
	for previousGuardPos := range visited {
		if previousGuardPos != guardStart {
			wg.Add(1)
			go func() {
				tryPossibleLoops(previousGuardPos)
				wg.Done()
			}()
		}
	}

	wg.Wait()
	fmt.Println("found", len(loopObstacles), "obstacles that cause loops")
}
