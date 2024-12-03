package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var (
	mul_pattern        = regexp.MustCompile(`(?:mul)\([0-9]{1,3}\,[0-9]{1,3}\)`)
	mul_switch_pattern = regexp.MustCompile(`(?:mul)\([0-9]{1,3}\,[0-9]{1,3}\)|do\(\)|don\'t\(\)`)
	num_pattern        = regexp.MustCompile(`([0-9]{1,3})`)
	total              int64
	switched_total     int64
)

func find_and_run_instructions(pattern *regexp.Regexp, data []byte, total *int64) {
	switched := true
	for _, instruction := range pattern.FindAll(data, 42069) {
		if string(instruction) == `do()` {
			switched = true
			continue
		} else if string(instruction) == `don't()` {
			switched = false
			continue
		}

		if !switched {
			continue
		}

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

		fmt.Printf("%d * %d = %d\n", n1, n2, n1*n2)
		*total += n1 * n2
	}
}

func main() {
	dat, err := os.ReadFile("input")
	if err != nil {
		panic(err)
	}

	find_and_run_instructions(mul_pattern, dat, &total)
	find_and_run_instructions(mul_switch_pattern, dat, &switched_total)

	fmt.Println(total, switched_total)
}
