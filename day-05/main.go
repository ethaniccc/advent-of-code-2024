package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	pageOrderRules = make(map[int64][]Ruleset)
	orders         = [][]int64{}
)

type Ruleset struct {
	Before, After int64
}

func (r Ruleset) Validate(order []int64) (int, bool) {
	beforeIndex, afterIndex := -1, -1
	for index, page := range order {
		if page == r.After {
			afterIndex = index
		} else if page == r.Before {
			beforeIndex = index
		}

		if beforeIndex != -1 && afterIndex != -1 {
			break
		}
	}

	return beforeIndex, (beforeIndex == -1 && afterIndex == -1) || (beforeIndex < afterIndex)
}

func parseInput() {
	dat, err := os.ReadFile("input")
	if err != nil {
		panic(err)
	}

	// The line below has caused me hours of suffering and pain. I didn't parse all the rules. This can't be real...
	for _, match := range regexp.MustCompile(`\d{1,2}\|\d{1,2}`).FindAll(dat, 10_000) {
		split := strings.Split(string(match), "|")
		printedBeforePg, _ := strconv.ParseInt(split[0], 10, 64)
		printPage, _ := strconv.ParseInt(split[1], 10, 64)

		if _, ok := pageOrderRules[printPage]; !ok {
			pageOrderRules[printPage] = []Ruleset{}
		}

		pageOrderRules[printPage] = append(pageOrderRules[printPage], Ruleset{Before: printedBeforePg, After: printPage})
	}

	for _, line := range regexp.MustCompile(`(\d{1,2},){1,200}\d{1,2}`).FindAll(dat, 1_000) {
		split := strings.Split(string(line), ",")
		order := make([]int64, len(split))

		for index, v := range split {
			page, _ := strconv.ParseInt(v, 10, 64)
			order[index] = page
		}
		orders = append(orders, order)
	}
}

func validateOrder(order []int64) bool {
	for _, page := range order {
		for _, rule := range pageOrderRules[page] {
			if _, ok := rule.Validate(order); !ok {
				return false
			}
		}
	}

	return true
}

func sortBadOrder(order []int64) {
	for index, page := range order {
		for _, rule := range pageOrderRules[page] {
			if beforeIndex, ok := rule.Validate(order); !ok {
				old := order[beforeIndex]
				order[beforeIndex] = page
				order[index] = old
				sortBadOrder(order)
				return
			}
		}
	}
}

func main() {
	parseInput()
	fmt.Println(len(orders), "orders")

	var validOrderSum, invalidOrderSum int64
	for _, order := range orders {
		if len(order)%2 == 0 {
			panic("will not be able to get order because of even length")
		}

		if validateOrder(order) {
			validOrderSum += order[len(order)/2]
		} else {
			fmt.Println("old:", order)
			sortBadOrder(order)
			if !validateOrder(order) {
				panic("nice one loser")
			}

			invalidOrderSum += order[len(order)/2]
			fmt.Println("new:", order)
		}
	}

	fmt.Println("valid order sum:", validOrderSum)
	fmt.Println("invalid order sum:", invalidOrderSum) // work wtf????
}
