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

func findBlockOfID(block []int64, id int64) (int, int) {
	endIndex := -1
	for index := len(block) - 1; index >= 0; index-- {
		if block[index] == id && endIndex == -1 {
			endIndex = index
		} else if block[index] != id && endIndex != -1 {
			return index, endIndex
		}
	}
	panic("block not found - this should never happen")
}

func findFreeBlockOfSize(block []int64, size int) (int, int, bool) {
	startIndex := -1
	for index, num := range block {
		if num == -1 && startIndex == -1 {
			startIndex = index
		} else if num != -1 && startIndex != -1 {
			if index-startIndex < size {
				startIndex = -1
			} else {
				return startIndex, index, true
			}
		}
	}
	return -1, -1, false
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

	var sum int64
	for index, num := range p1block {
		if num == -1 {
			break
		}
		sum += int64(index) * num
	}
	fmt.Println("part1 checksum:", sum)
}

func part2() {
	p2block := make([]int64, len(block))
	copy(p2block, block)
	fmt.Println(p2block)

	for id := int64(9); id >= 0; id-- {
		startOfBlock, endOfBlock := findBlockOfID(p2block, id)
		startOfMove, endOfMove, foundMove := findFreeBlockOfSize(p2block, (endOfBlock-startOfBlock)+1)
		if !foundMove || startOfMove > startOfBlock {
			continue
		}

		fmt.Println(startOfMove, endOfMove, startOfBlock, endOfBlock)
		oldIndex := startOfBlock
		for i := startOfMove; i <= endOfMove; i++ {
			p2block[i] = p2block[oldIndex]
			p2block[oldIndex] = -1
			oldIndex++
			fmt.Println(oldIndex, i)
		}
	}
	fmt.Println(p2block)
}

func main() {
	parseInput()
	part1()
	part2()
}
