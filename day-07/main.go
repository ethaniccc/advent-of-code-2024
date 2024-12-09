package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
)

type EquationNode struct {
	Num  uint64
	Next *EquationNode
}

func NewEquationNode(nums []uint64) *EquationNode {
	n := &EquationNode{Num: nums[0]}
	currentNode := n
	for index := 1; index < len(nums); index++ {
		currentNode.Next = &EquationNode{Num: nums[index]}
		currentNode = currentNode.Next
	}

	return n
}

func (n *EquationNode) Results() []uint64 {
	results := []uint64{n.Num}
	concatedResults := []uint64{}
	nextNode := n.Next
	//previousNodeNum := n.Num

	for nextNode != nil {
		newResults := []uint64{}
		for _, num := range results {
			//pr1, pr2 := num/previousNodeNum, num-previousNodeNum
			newResults = append(
				newResults,
				num*nextNode.Num,
				num+nextNode.Num,
				concate(num, nextNode.Num),
				//concate(pr1*previousNodeNum, num),
				//concate(pr2+previousNodeNum, num),
				//concate(pr1+previousNodeNum, num),
				//concate(pr2*previousNodeNum, num),
			)
		}

		results = results[:0]
		results = append(results, newResults...)
		//previousNodeNum = nextNode.Num
		nextNode = nextNode.Next
	}

	return append(results, concatedResults...)
}

type Equation struct {
	Nums           []uint64
	ExpectedResult uint64
}

var equations []Equation

func parseInput() {
	fName := "input"
	if os.Getenv("TEST") != "" {
		fName = "test-input"
	}

	dat, err := os.ReadFile(fName)
	if err != nil {
		panic(err)
	}

	results := regexp.MustCompile(`\d{1,99}\:`).FindAll(dat, 1_000_000)
	eqs := regexp.MustCompile(`(\d{1,3}\s){1,100}`).FindAll(dat, 1_000_000)

	equations = make([]Equation, len(eqs))
	for index, line := range eqs {
		l := strings.ReplaceAll(string(line), "\r", "")
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

		equations[index] = Equation{
			Nums:           nums,
			ExpectedResult: expectedOutcome,
		}
	}
}

func concate(n1, n2 uint64) uint64 {
	var base uint64
	for base = 1; base < n2; base *= 10 {
	}

	return (n1 * base) + n2
}

func main() {
	parseInput()

	debug.SetMemoryLimit(4 * 1024 * 1024)
	debug.SetGCPercent(-1)

	validEquations, totalValidResult := 0, uint64(0)
	eqs := 0
	for _, eq := range equations {
		eqs++
		results := NewEquationNode(eq.Nums).Results()
		if slices.Contains(results, eq.ExpectedResult) {
			validEquations++
			totalValidResult += eq.ExpectedResult
			fmt.Println("equation", eqs, "is valid")
		} else {
			fmt.Println("equation", eqs, "is not valid")
		}
	}

	// 482740316327806 is not answer
	fmt.Println(validEquations, "valid equations totaling", totalValidResult)
}
