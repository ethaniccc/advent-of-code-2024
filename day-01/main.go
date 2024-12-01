package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var (
	left_list  = make([]int32, 1000)
	right_list = make([]int32, 1000)
)

func parse_input() {
	// Read the data from the input file.
	dat, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	// Reach each line and parse the left and right list.
	for index, line := range strings.Split(string(dat), "\n") {
		// Split the line from the whitespace inbetween the left and right lists.
		split := strings.Fields(line)
		left_num, _ := strconv.ParseUint(split[0], 10, 64)
		right_num, _ := strconv.ParseUint(split[1], 10, 64)

		// Put each parsed integer into the left & right list.
		left_list[index] = int32(left_num)
		right_list[index] = int32(right_num)
	}
}

func main() {
	// Parse the input and fill the left/right lists.
	parse_input()

	// Sort both the slices using the slices package.
	slices.Sort(left_list)
	slices.Sort(right_list)

	var total_distance int32
	for index := 0; index < 1_000; index++ {
		left_num := left_list[index]
		right_num := right_list[index]

		distance := left_num - right_num
		if distance < 0 {
			distance *= -1
		}

		total_distance += distance
	}

	fmt.Println(total_distance)
}
