package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type EquationNode struct {
	Num       uint64
	Expecting uint64
	Next      *EquationNode
}

func NewEquationNode(nums []uint64, expecting uint64) *EquationNode {
	n := &EquationNode{Num: nums[0], Expecting: expecting}
	currentNode := n
	for index := 1; index < len(nums); index++ {
		currentNode.Next = &EquationNode{Num: nums[index]}
		currentNode = currentNode.Next
	}

	return n
}

func (n *EquationNode) Test() bool {
	results := []uint64{n.Num}
	for nextNode := n.Next; nextNode != nil; nextNode = nextNode.Next {
		newResults := []uint64{}
		for _, num := range results {
			newResults = append(
				newResults,
				num*nextNode.Num,
				num+nextNode.Num,
				concate(num, nextNode.Num),
			)
		}

		results = results[:0]
		results = append(results, newResults...)
	}

	for _, result := range results {
		if result == n.Expecting {
			return true
		}
	}
	return false
}

var equations []*EquationNode

func parseInput() {
	fName := "input"
	if os.Getenv("TEST") != "" {
		fName = "test-input"
	}

	dat, err := os.ReadFile(fName)
	if err != nil {
		panic(err)
	}

	newline := "\n"
	if runtime.GOOS == "windows" {
		newline = "\r"
	}

	results := regexp.MustCompile(`\d{1,99}\:`).FindAll(dat, 1_000_000)
	eqs := regexp.MustCompile(`(\d{1,3}\s){1,100}`).FindAll(dat, 1_000_000)

	equations = make([]*EquationNode, len(eqs))
	for index, line := range eqs {
		l := strings.ReplaceAll(string(line), newline, "")
		split := strings.Split(l, " ")
		nums := make([]uint64, len(split))

		for i2, v := range split {
			if v == "" {
				continue
			}

			num, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				panic(err)
			}

			if num == 0 {
				panic("num is 0")
			}

			nums[i2] = num
		}

		expectedOutcome, err := strconv.ParseUint(strings.ReplaceAll(string(results[index]), ":", ""), 10, 64)
		if err != nil {
			panic(err)
		}

		equations[index] = NewEquationNode(nums, expectedOutcome)
	}
}

func concate(n1, n2 uint64) uint64 {
	var base uint64
	for base = 1; base <= n2; base *= 10 {
		if base == 0 {
			panic("uint64 overflow when concating")
		}
	}

	return (n1 * base) + n2
}

func main() {
	parseInput()

	var wg sync.WaitGroup
	validEquations, totalValidResult := 0, uint64(0)
	for _, eq := range equations {
		wg.Add(1)
		go func(eq *EquationNode) {
			if eq.Test() {
				validEquations++
				totalValidResult += eq.Expecting
			}
			wg.Done()
		}(eq)
	}

	// 482740316327806 is not answer
	// 482739998262825 is not answer
	wg.Wait()
	fmt.Println(validEquations, "valid equations totaling", totalValidResult)
}
