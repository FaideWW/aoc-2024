package main

import (
	"fmt"
	"os"

	io "github.com/faideww/aoc-2024/lib"
)

type Position struct {
	x int
	y int
}

type Board struct {
	width    int
	height   int
	antennae map[rune][]Position
}

func main() {
	input := io.ReadInputFile(os.Args[1])
	board := parseBoard(input)

	result := countAntinodes(board)
	fmt.Printf("unique antinodes (part1): %d\n", result)

	result2 := countAllAntinodes(board)
	fmt.Printf("unique antinodes (part2): %d\n", result2)
}

func parseBoard(input string) Board {
	antennae := make(map[rune][]Position)
	lines := io.TrimAndSplit(input)
	for y, line := range lines {
		for x, char := range line {
			if char == '.' {
				continue
			}
			if _, ok := antennae[char]; !ok {
				antennae[char] = []Position{}
			}
			antennae[char] = append(antennae[char], Position{x, y})
		}
	}

	return Board{
		width:    len(lines[0]),
		height:   len(lines),
		antennae: antennae,
	}
}

func countAntinodes(board Board) int {
	allAntinodes := make(map[Position]bool)
	for freq := range board.antennae {
		antinodes := findAntinodes(board, freq)

		for pos := range antinodes {
			allAntinodes[pos] = true
		}
	}

	return len(allAntinodes)
}

func findAntinodes(board Board, freq rune) map[Position]bool {
	// iterate through all pairs of antennae and compute the antinodes
	antennae := board.antennae[freq]

	antinodes := make(map[Position]bool)

	for i := 0; i < len(antennae); i++ {
		node1 := antennae[i]
		for j := i + 1; j < len(antennae); j++ {
			node2 := antennae[j]
			dx := node2.x - node1.x
			dy := node2.y - node1.y

			antinode1 := Position{node2.x + dx, node2.y + dy}

			if antinode1.x >= 0 && antinode1.x < board.width && antinode1.y >= 0 && antinode1.y < board.height {
				antinodes[antinode1] = true
			}

			antinode2 := Position{node1.x - dx, node1.y - dy}

			if antinode2.x >= 0 && antinode2.x < board.width && antinode2.y >= 0 && antinode2.y < board.height {
				antinodes[antinode2] = true
			}
		}
	}

	return antinodes
}

func countAllAntinodes(board Board) int {
	allAntinodes := make(map[Position]bool)
	for freq := range board.antennae {
		antinodes := findAllAntinodes(board, freq)

		for pos := range antinodes {
			allAntinodes[pos] = true
		}
	}

	return len(allAntinodes)
}

func findAllAntinodes(board Board, freq rune) map[Position]bool {
	antennae := board.antennae[freq]

	antinodes := make(map[Position]bool)

	for i := 0; i < len(antennae); i++ {
		node1 := antennae[i]
		for j := i + 1; j < len(antennae); j++ {
			node2 := antennae[j]
			dx := node2.x - node1.x
			dy := node2.y - node1.y

			// count antinodes in 1 direction
			currentPos := Position{node2.x, node2.y}
			for currentPos.x >= 0 && currentPos.x < board.width && currentPos.y >= 0 && currentPos.y < board.height {
				antinodes[currentPos] = true
				currentPos.x += dx
				currentPos.y += dy
			}

			// now count them in the other direction
			currentPos = Position{node1.x, node1.y}
			for currentPos.x >= 0 && currentPos.x < board.width && currentPos.y >= 0 && currentPos.y < board.height {
				antinodes[currentPos] = true
				currentPos.x -= dx
				currentPos.y -= dy
			}
		}
	}

	return antinodes
}
