package main

import (
	"fmt"
	"os"
	"strconv"
)

var block []int64

func parseInput() {
	dat, err := os.ReadFile("input")
	if err != nil {
		panic(err)
	}

	var currentBlock int64
	for index := 0; index < len(dat)-1; index++ {
		// If the index divides evenly into 2, then this is a repeating block
		value, err := strconv.ParseInt(string(dat[index]), 10, 64)
		if err != nil {
			panic(err)
		}

		for v := value; v > 0; v-- {
			if index%2 == 0 {
				block = append(block, currentBlock)
			} else {
				block = append(block, -1)
			}
		}

		if index%2 == 0 {
			currentBlock++
		}
	}
}

func findFreeSpaceFromStart(block []int64) int {
	for index := 0; index < len(block); index++ {
		if block[index] == -1 {
			return index
		}
	}
	panic("no free space found from start")
}

func findFreeSpaceFromEnd(block []int64) int {
	oldIndex := len(block) - 1
	for index := len(block) - 1; index >= 0; index-- {
		if block[index] == -1 {
			oldIndex = index
		} else {
			return oldIndex
		}
	}
	panic("no free space found from end")
}

func findNumFromEnd(block []int64) int {
	for index := len(block) - 1; index >= 0; index-- {
		if block[index] != -1 {
			return index
		}
	}
	panic("no number found from end")
}

func findBlockOfID(block []int64, id int64) (int, int, bool) {
	endIndex := -1
	for index := len(block) - 1; index >= 0; index-- {
		if block[index] == id && endIndex == -1 {
			endIndex = index
		} else if block[index] != id && endIndex != -1 {
			return index + 1, endIndex, true
		}
	}

	if endIndex != -1 {
		return 0, endIndex, true
	}
	return -1, -1, false
}

func findFreeBlockOfSize(block []int64, size int) (int, int, bool) {
	startIndex := -1
	for index, num := range block {
		if num == -1 && startIndex == -1 {
			if size == 0 {
				return index, index, true
			}
			startIndex = index
		} else if num != -1 && startIndex != -1 {
			if (index-1)-startIndex < size {
				startIndex = -1
			} else {
				return startIndex, index - 1, true
			}
		}
	}

	if startIndex != -1 {
		return startIndex, len(block) - 1, true
	}
	return -1, -1, false
}

func blockChecksum(block []int64) int64 {
	var sum int64
	for index, num := range block {
		if num == -1 {
			continue
		}
		sum += int64(index) * num
	}

	return sum
}

func part1() {
	p1block := make([]int64, len(block))
	copy(p1block, block)

	for {
		nextFreeSpot := findFreeSpaceFromStart(p1block)
		if spaceFromEnd := findFreeSpaceFromEnd(p1block); nextFreeSpot == spaceFromEnd {
			break
		}

		nextNumIndex := findNumFromEnd(p1block)
		p1block[nextFreeSpot] = p1block[nextNumIndex]
		p1block[nextNumIndex] = -1
	}

	fmt.Println("part1 checksum:", blockChecksum(p1block))
}

func part2() {
	p2block := make([]int64, len(block))
	copy(p2block, block)

	for id := int64(10000); id >= 0; id-- {
		startOfBlock, endOfBlock, blockFound := findBlockOfID(p2block, id)
		if !blockFound {
			fmt.Println("block", id, "not found")
			continue
		}

		startOfMove, _, foundMove := findFreeBlockOfSize(p2block, endOfBlock-startOfBlock)
		if !foundMove || startOfMove > startOfBlock {
			continue
		}

		moveIndex := startOfMove
		for i := startOfBlock; i <= endOfBlock; i++ {
			p2block[moveIndex] = p2block[i]
			p2block[i] = -1
			moveIndex++
		}
	}
	fmt.Println("part 2 checksum:", blockChecksum(p2block))
}

func main() {
	parseInput()
	part1()
	part2()
}
