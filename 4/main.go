package main

import (
	"fmt"
	"os"

	io "github.com/faideww/aoc-2024/lib"
)

func main() {
	input := io.ReadInputFile(os.Args[1])
	lines := io.TrimAndSplit(input)
	result := doXmasWordSearch(lines)

	fmt.Printf("xmas search results: %d\n", result)

	result2 := doXDashMasWordSearch(lines)

	fmt.Printf("x-mas search results: %d\n", result2)
}

func doXmasWordSearch(lines []string) int {
	sum := 0
	// search horizontally
	for y, line := range lines {
		for x, char := range line {
			if char == 'X' {
				if x+3 < len(line) && line[x+1] == 'M' && line[x+2] == 'A' && line[x+3] == 'S' {
					// horizontal forward line
					sum++
				}

				if y+3 < len(lines) && lines[y+1][x] == 'M' && lines[y+2][x] == 'A' && lines[y+3][x] == 'S' {
					// vertical forward line
					sum++
				}

				if x+3 < len(line) && y+3 < len(lines) && lines[y+1][x+1] == 'M' && lines[y+2][x+2] == 'A' && lines[y+3][x+3] == 'S' {
					// diagonal down-forward line
					sum++
				}

				if x+3 < len(line) && y >= 3 && lines[y-1][x+1] == 'M' && lines[y-2][x+2] == 'A' && lines[y-3][x+3] == 'S' {
					// diagonal up-forward line
					sum++
				}

			} else if char == 'S' {
				if x+3 < len(line) && line[x+1] == 'A' && line[x+2] == 'M' && line[x+3] == 'X' {
					// horizontal reverse line
					sum++
				}

				if y+3 < len(lines) && lines[y+1][x] == 'A' && lines[y+2][x] == 'M' && lines[y+3][x] == 'X' {
					// vertical reverse line
					sum++
				}

				if x+3 < len(line) && y+3 < len(lines) && lines[y+1][x+1] == 'A' && lines[y+2][x+2] == 'M' && lines[y+3][x+3] == 'X' {
					// diagonal down-reverse line
					sum++
				}

				if x+3 < len(line) && y >= 3 && lines[y-1][x+1] == 'A' && lines[y-2][x+2] == 'M' && lines[y-3][x+3] == 'X' {
					// diagonal up-reverse line
					sum++
				}
			}
		}
	}

	return sum
}

func doXDashMasWordSearch(lines []string) int {
	sum := 0
	// search horizontally
	for y := 1; y < len(lines)-1; y++ {
		line := lines[y]
		for x := 1; x < len(line)-1; x++ {
			char := line[x]
			if char == 'A' {
				downRightOk := false
				upRightOk := false
				if (lines[y-1][x-1] == 'M' && lines[y+1][x+1] == 'S') || (lines[y-1][x-1] == 'S' && lines[y+1][x+1] == 'M') {
					downRightOk = true
				}
				if (lines[y+1][x-1] == 'M' && lines[y-1][x+1] == 'S') || (lines[y+1][x-1] == 'S' && lines[y-1][x+1] == 'M') {
					upRightOk = true
				}

				if downRightOk && upRightOk {
					sum++
				}
			}
		}
	}

	return sum
}
