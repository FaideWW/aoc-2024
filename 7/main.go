package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2024/lib"
)

func main() {
	input := io.ReadInputFile(os.Args[1])
	lines := io.TrimAndSplit(input)

	result := findValidEquations(lines)
	fmt.Printf("valid equation sum: %d\n", result)

}

func findValidEquations(lines []string) (sum int) {
	for _, line := range lines {
		fields := strings.Fields(line)
		resultStr := fields[0]
		target, _ := strconv.Atoi(resultStr[:len(resultStr)-1])
		operands := make([]int, len(fields)-1)
		for i, opStr := range fields[1:] {
			operands[i], _ = strconv.Atoi(opStr)
		}

		if searchValidEquation(target, operands, operands[0], 1) {
			sum += target
		}
	}
	return sum
}

func searchValidEquation(target int, operands []int, currentValue int, idx int) bool {
	// recursively reduce the equation by applying one of the four operators until we find an equality

	if idx == len(operands) {
		return target == currentValue
	}

	// addition
	if searchValidEquation(target, operands, currentValue+operands[idx], idx+1) {
		return true
	}

	// multiplication
	if searchValidEquation(target, operands, currentValue*operands[idx], idx+1) {
		return true
	}

	// concatenation (pt 2)
	nextValue, _ := strconv.Atoi(fmt.Sprintf("%d%d", currentValue, operands[idx]))
	if searchValidEquation(target, operands, nextValue, idx+1) {
		return true
	}
	return false
}
