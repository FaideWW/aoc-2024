package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2024/lib"
)

func main() {
	input := io.ReadInputFile(os.Args[1])
	lines := io.TrimAndSplit(input)

	sum := sumSortedPairs(lines)
	fmt.Printf("sum: %d\n", sum)

	similarity := computeSimilarityScore(lines)
	fmt.Printf("similarity: %d\n", similarity)
}

func sumSortedPairs(lines []string) int {
	list1 := make([]int, len(lines))
	list2 := make([]int, len(lines))

	for i, line := range lines {
		values := strings.Fields(line)
		list1[i], _ = strconv.Atoi(values[0])
		list2[i], _ = strconv.Atoi(values[1])
	}

	slices.Sort(list1)
	slices.Sort(list2)

	sum := 0
	for i := range list1 {
		v1 := list1[i]
		v2 := list2[i]

		if v2 > v1 {
			sum += v2 - v1
		} else {
			sum += v1 - v2
		}
	}

	return sum
}

func computeSimilarityScore(lines []string) int {
	list1 := make([]int, len(lines))

	seen := make(map[int]int)

	for i, line := range lines {
		values := strings.Fields(line)
		list1[i], _ = strconv.Atoi(values[0])
		val2, _ := strconv.Atoi(values[1])

		if _, ok := seen[val2]; !ok {
			seen[val2] = 0
		}

		seen[val2]++
	}

	sum := 0
	for _, val := range list1 {
		if count, ok := seen[val]; ok {
			sum += val * count
		}
	}
	return sum
}
