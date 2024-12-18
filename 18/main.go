package main

import (
	"fmt"
	"os"
	"strconv"

	io "github.com/faideww/aoc-2024/lib"
)

type Position struct{ x, y int }

const (
	EMPTY_TILE = iota
	CORRUPTED_TILE
)

type Board struct {
	width, height int
	bytes         []Position
	tiles         map[Position]struct{}
}

func main() {
	input := io.ReadInputFile(os.Args[1])
	bytesSimulated, _ := strconv.Atoi(os.Args[2])
	arenaSize, _ := strconv.Atoi(os.Args[3])

	board := parseBoard(input, arenaSize)
	for i := 0; i < bytesSimulated; i++ {
		simulateByte(&board, i)
	}

	result, hasPath := findShortestPath(board)
	fmt.Printf("shortest path has %d unique tiles\n", result)

	nextI := bytesSimulated - 1
	for hasPath && nextI < len(board.bytes)-1 {
		nextI++
		simulateByte(&board, nextI)
		_, hasPath = findShortestPath(board)
	}

	if hasPath {
		fmt.Printf("no blocking bytes found??\n")
	} else {
		fmt.Printf("first blocking byte: %d (%v)\n", nextI, board.bytes[nextI])
	}
}

func parseBoard(input string, arenaSize int) Board {
	lines := io.TrimAndSplit(input)
	bytes := make([]Position, len(lines))
	for i := range lines {
		components := io.TrimAndSplitBy(lines[i], ",")
		x, _ := strconv.Atoi(components[0])
		y, _ := strconv.Atoi(components[1])

		bytes[i] = Position{x, y}
	}

	return Board{
		width:  arenaSize,
		height: arenaSize,
		bytes:  bytes,
		tiles:  make(map[Position]struct{}),
	}
}

func simulateByte(board *Board, byteIndex int) {
	board.tiles[board.bytes[byteIndex]] = struct{}{}
}

func findShortestPath(board Board) (int, bool) {
	// Pure, straightforward BFS

	start := Position{0, 0}
	goal := Position{board.width - 1, board.height - 1}

	frontier := make([]Position, 0)
	explored := make(map[Position]struct{})

	cameFrom := make(map[Position]Position)

	frontier = append(frontier, start)
	explored[start] = struct{}{}

	foundGoal := false

	for len(frontier) > 0 {
		current := frontier[0]
		frontier = frontier[1:]

		if current == goal {
			foundGoal = true
			break
		}

		neighbors := findNeighbors(board, current)
		for _, neighbor := range neighbors {
			if _, seen := explored[neighbor]; !seen {
				explored[neighbor] = struct{}{}
				cameFrom[neighbor] = current
				frontier = append(frontier, neighbor)
			}
		}
	}

	if !foundGoal {
		return -1, false
	}

	pathMap := make(map[Position]struct{})

	current := goal
	for current != start {
		pathMap[current] = struct{}{}
		current = cameFrom[current]
	}

	// printBoard(board, pathMap)

	return len(pathMap), true
}

func findNeighbors(board Board, current Position) []Position {
	neighbors := []Position{}

	up := Position{current.x, current.y - 1}
	right := Position{current.x + 1, current.y}
	down := Position{current.x, current.y + 1}
	left := Position{current.x - 1, current.y}

	_, upCorrupted := board.tiles[up]
	_, rightCorrupted := board.tiles[right]
	_, downCorrupted := board.tiles[down]
	_, leftCorrupted := board.tiles[left]

	if current.x > 0 && !leftCorrupted {
		neighbors = append(neighbors, left)
	}
	if current.x < board.width-1 && !rightCorrupted {
		neighbors = append(neighbors, right)
	}
	if current.y > 0 && !upCorrupted {
		neighbors = append(neighbors, up)
	}
	if current.y < board.height-1 && !downCorrupted {
		neighbors = append(neighbors, down)
	}
	return neighbors
}

func printBoard(board Board, explored map[Position]struct{}) {
	for y := 0; y < board.height; y++ {
		for x := 0; x < board.width; x++ {
			p := Position{x, y}
			if _, corrupted := board.tiles[p]; corrupted {
				fmt.Print("#")
			} else if _, visited := explored[p]; visited {
				fmt.Print("O")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
