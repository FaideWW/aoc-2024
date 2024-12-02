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

	sum1 := computeSafeReports(lines)
	fmt.Printf("safe reports: %d\n", sum1)

	sum2 := computeSafeReportsWithTolerance(lines)
	fmt.Printf("safe reports with tolerance: %d\n", sum2)

}

func computeSafeReports(lines []string) int {
	sum := 0
	for _, line := range lines {
		levels := make([]int, 0)
		for _, sVar := range strings.Fields(line) {
			iVar, _ := strconv.Atoi(sVar)
			levels = append(levels, iVar)
		}

		if isReportSafe(levels) {
			sum++
		}
	}

	return sum
}

func computeSafeReportsWithTolerance(lines []string) int {
	sum := 0
	for _, line := range lines {
		levels := make([]int, 0)
		for _, sVar := range strings.Fields(line) {
			iVar, _ := strconv.Atoi(sVar)
			levels = append(levels, iVar)
		}

		if isReportSafe(levels) {
			sum++
		} else {
			for j := 0; j < len(levels); j++ {
				// annoying slice manipulation (append modifies the underlying
				// array for performance reasons)
				levelsWithoutBadIndex := slices.Delete(slices.Clone(levels), j, j+1)
				if isReportSafe(levelsWithoutBadIndex) {
					sum++
					break
				}
			}
		}
	}

	return sum
}

func isReportSafe(levels []int) bool {
	var inc bool
	if levels[0] < levels[1] {
		inc = true
	} else if levels[0] > levels[1] {
		inc = false
	} else {
		return false
	}

	for i := 1; i < len(levels); i++ {
		if (inc && levels[i-1] >= levels[i]) || (!inc && levels[i-1] <= levels[i]) {
			return false
		}

		if levels[i-1]-levels[i] > 3 || levels[i-1]-levels[i] < -3 {
			return false
		}
	}

	return true
}
