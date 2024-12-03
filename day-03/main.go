package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var (
	mul_pattern = regexp.MustCompile(`(?:mul)\([0-9]{1,3}\,[0-9]{1,3}\)`)
	// TODO: complete this because regex will be the end of me...
	mul_with_switch_pattern = regexp.MustCompile(`(?:mul)\([0-9]{1,3}\,[0-9]{1,3}\)`)
	num_pattern             = regexp.MustCompile(`([0-9]{1,3})`)
	total                   int64
	switched_total          int64
)

func find_and_run_instructions(pattern *regexp.Regexp, data []byte, total *int64) {
	for _, instruction := range pattern.FindAll(data, 42069) {
		var (
			n1, n2 int64
			err    error
		)
		nums := num_pattern.FindAllString(string(instruction), 2)
		if n1, err = strconv.ParseInt(nums[0], 10, 64); err != nil {
			panic(err)
		}
		if n2, err = strconv.ParseInt(nums[1], 10, 64); err != nil {
			panic(err)
		}

		*total += n1 * n2
	}
}

func main() {
	dat, err := os.ReadFile("input")
	if err != nil {
		panic(err)
	}

	sample := "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))"
	fmt.Println(mul_with_switch_pattern.FindAllString(sample, 4))

	find_and_run_instructions(mul_pattern, dat, &total)
	find_and_run_instructions(mul_with_switch_pattern, dat, &switched_total)
	fmt.Println(total, switched_total, switched_total < total)
}
