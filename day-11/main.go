package main

import (
	"fmt"
)

type blinkArgs [2]uint64

// SO: first I attempted to use just slices but that was taking all my memory for part 2, so I
// moved on to trying some node stuff, and it saved my memory a bunch - but the runtime was
// still taking too long. couldn't figure out how I could possibly solve the runtime issue until
// someone on the reddit mentioned caching. I instantly changed everything into a recursive function
// so I could cache the resuts, which made the runtime shrink into a few milliseonds.
// It is crazy how AOC and it's puzzles are very humbling and teach me that I should just focus
// on what is needed as the result (which in this case, is not the actual stone numbers, but just
// the NUMBBER OF STONES), and that caching is FUCKING AWESOME!!!
var cache = make(map[blinkArgs]uint64)

func split(num uint64) (
	uint64,
	uint64,
	bool,
) {
	var power, digits uint64 = 1, 0
	for power <= num {
		power *= 10
		digits++
	}
	power /= 10
	if digits%2 != 0 {
		return 0, 0, false
	}

	splitConstant := pow(10, digits/2)
	leftSplit := num / splitConstant
	rightSplit := num - (leftSplit * splitConstant)

	return leftSplit, rightSplit, true
}

func pow(
	n,
	base uint64,
) uint64 {
	var v uint64 = 1
	for base > 0 {
		v *= n
		base--
	}
	return v
}

func blink(
	stone,
	times uint64,
) uint64 {
	if times == 0 {
		return 1
	} else if result, isCached := cache[blinkArgs{stone, times}]; isCached {
		return result
	}

	times--
	if stone == 0 {
		cache[blinkArgs{1, times}] = blink(1, times)
		return cache[blinkArgs{1, times}]
	} else if left, right, ableToSplit := split(stone); ableToSplit {
		cache[blinkArgs{left, times}] = blink(left, times)
		cache[blinkArgs{right, times}] = blink(right, times)
		return cache[blinkArgs{left, times}] + cache[blinkArgs{right, times}]
	} else {
		cache[blinkArgs{stone * 2024, times}] = blink(stone*2024, times)
		return cache[blinkArgs{stone * 2024, times}]
	}
}

func main() {
	var totalStones uint64
	stones := []uint64{3, 386358, 86195, 85, 1267, 3752457, 0, 741}
	for _, stone := range stones {
		totalStones += blink(stone, 75)
	}

	fmt.Println(totalStones)
}
