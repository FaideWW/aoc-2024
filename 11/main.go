package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	io "github.com/faideww/aoc-2024/lib"
)

func main() {
	input := io.ReadInputFile(os.Args[1])
	stones := parseStones(input)

	stoneCount := blink(stones, 25)
	fmt.Printf("stones: %d\n", stoneCount)

	stones = parseStones(input)

	stoneCount2 := blink2(stones, 75)
	fmt.Printf("stones: %d\n", stoneCount2)
}

func parseStones(input string) []int {
	stoneStrs := strings.Fields(input)

	stones := make([]int, len(stoneStrs))
	for i, str := range stoneStrs {
		value, _ := strconv.Atoi(str)
		stones[i] = value
	}
	return stones
}

func blink(stones []int, n int) int {
	funcBegin := time.Now()

	for i := 0; i < n; i++ {
		loopBegin := time.Now()
		for j := 0; j < len(stones); j++ {
			stone := stones[j]
			magnitude := ord(stone)
			if stone == 0 {
				stones[j] = 1
			} else if magnitude%2 == 0 {
				halfMag := tenToThe(magnitude / 2)
				stone1 := stone / halfMag
				stone2 := stone % halfMag
				stones = append(stones[:j+1], stones[j:]...)
				stones[j] = stone1
				stones[j+1] = stone2
				j++
			} else {
				stones[j] = stone * 2024
			}
		}
		loopTime := time.Since(loopBegin)
		fmt.Printf("blink %d took %s\n", i+1, loopTime)
	}
	funcTime := time.Since(funcBegin)
	fmt.Printf("took %s\n", funcTime)
	return len(stones)
}

func blink2(stones []int, n int) int {
	// Hypothesis: the rules will naturally result in a LOT of repeating values
	// (since we always split larger numbers into smaller numbers until they're
	// single digits). So instead of stepping through every single stone, we step
	// through "classes" of stones (eg. compute the result of blinking a 1, 2, 3,
	// etc.) while keeping track of how many of each class exist. Then we add
	// that many of the results of the blink to the next iteration.
	funcBegin := time.Now()
	stoneMap := make(map[int]int)
	for _, stone := range stones {
		if _, ok := stoneMap[stone]; !ok {
			stoneMap[stone] = 0
		}
		stoneMap[stone]++
	}

	for i := 0; i < n; i++ {
		loopBegin := time.Now()
		nextStoneMap := make(map[int]int)
		for stone, count := range stoneMap {
			magnitude := ord(stone)
			if stone == 0 {
				addToMap(nextStoneMap, 1, count)
			} else if magnitude%2 == 0 {
				halfMag := tenToThe(magnitude / 2)
				stone1 := stone / halfMag
				stone2 := stone % halfMag

				addToMap(nextStoneMap, stone1, count)
				addToMap(nextStoneMap, stone2, count)
			} else {
				addToMap(nextStoneMap, stone*2024, count)
			}
		}

		stoneMap = nextStoneMap

		loopTime := time.Since(loopBegin)
		fmt.Printf("blink %d took %s\n", i+1, loopTime)
	}

	sum := 0
	for _, count := range stoneMap {
		sum += count
	}
	funcTime := time.Since(funcBegin)
	fmt.Printf("took %s\n", funcTime)
	return sum
}

func addToMap(m map[int]int, key, count int) {
	if _, ok := m[key]; !ok {
		m[key] = 0
	}
	m[key] += count
}

func ord(value int) int {
	order := 1
	for value >= 10 {
		value = value / 10
		order++
	}

	return order
}

func tenToThe(n int) int {
	v := 10
	for i := 1; i < n; i++ {
		v *= 10
	}
	return v
}
