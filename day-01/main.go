package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var (
	left_list  = make([]int64, 1000)
	right_list = make([]int64, 1000)
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
		left_num, _ := strconv.ParseInt(split[0], 10, 64)
		right_num, _ := strconv.ParseInt(split[1], 10, 64)

		// Put each parsed integer into the left & right list.
		left_list[index] = left_num
		right_list[index] = right_num
	}
}

func calculate_distance() int64 {
	var total_distance int64

	for index := 0; index < 1_000; index++ {
		left_num := left_list[index]
		right_num := right_list[index]

		distance := left_num - right_num
		if distance < 0 {
			distance *= -1
		}

		total_distance += distance
	}

	return total_distance
}

func calculate_similarity() int64 {
	var (
		apperances       = make(map[int64]uint16)
		similarity_score int64
	)

	// First iterate through the right list, to find the amount of times certain numbers appear
	// and put them on the apperances map.
	for _, num := range right_list {
		if _, ok := apperances[num]; !ok {
			apperances[num] = 1
			continue
		}

		apperances[num]++
	}

	// Now, iterate through the left list to see how many times a certain number appears on the
	// right list, and add to the similarity score the times appeared multiplied by the number itself.
	for _, num := range left_list {
		if times_appeared, ok := apperances[num]; ok {
			similarity_score += int64(times_appeared) * num
		}
	}

	return similarity_score
}

func main() {
	// Parse the input and fill the left/right lists.
	parse_input()

	// Sort both the slices using the slices package.
	slices.Sort(left_list)
	slices.Sort(right_list)

	fmt.Println("distance:", calculate_distance())
	fmt.Println("similarity:", calculate_similarity())
}
