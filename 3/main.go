package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2024/lib"
)

func main() {
	input := strings.TrimSpace(io.ReadInputFile(os.Args[1]))
	result := computeValidMuls(input)
	fmt.Printf("sum muls: %d\n", result)

	result2 := computeValidMulsStateful(input)
	fmt.Printf("sum muls with state: %d\n", result2)
}

func computeValidMuls(input string) int {
	r := regexp.MustCompile(`mul\((\d+),(\d+)\)`)

	sum := 0
	for _, match := range r.FindAllStringSubmatch(input, -1) {
		v1, _ := strconv.Atoi(match[1])
		v2, _ := strconv.Atoi(match[2])
		sum += v1 * v2
	}

	return sum
}

func computeValidMulsStateful(input string) int {

	r := regexp.MustCompile(`(mul|do|don't)\((?:(\d+),(\d+))?\)`)

	sum := 0
	on := true
	for _, match := range r.FindAllStringSubmatch(input, -1) {
		switch match[1] {
		case "mul":
			if on {
				v1, _ := strconv.Atoi(match[2])
				v2, _ := strconv.Atoi(match[3])
				sum += v1 * v2
			}
		case "do":
			on = true
		case "don't":
			on = false
		}
	}

	return sum
}
