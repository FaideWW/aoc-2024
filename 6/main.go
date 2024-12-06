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

type GuardState struct {
	pos    Position
	facing int // 0 = up, 1 = right, 2 = down, 3 = left
}

type Board struct {
	width      int
	height     int
	obstacles  map[Position]bool
	guardState GuardState
}

func main() {
	input := io.ReadInputFile(os.Args[1])
	board := parseBoard(input)
	uniqueTiles := traceGuardPath(board)
	loopsFound := findLoops(board, uniqueTiles)

	fmt.Printf("unique tiles crossed: %d\n", len(uniqueTiles))
	fmt.Printf("loops found: %d\n", loopsFound)
}

func parseBoard(input string) Board {
	lines := io.TrimAndSplit(input)
	height := len(lines)
	width := len(lines[0])
	obstacles := make(map[Position]bool)
	guardState := GuardState{facing: 0}
	for y, line := range lines {
		for x, char := range line {
			switch char {
			case '#':
				obstacles[Position{x, y}] = true
			case '^':
				guardState.pos = Position{x, y}
			}
		}
	}

	return Board{
		height:     height,
		width:      width,
		obstacles:  obstacles,
		guardState: guardState,
	}
}

func traceGuardPath(board Board) map[Position]bool {
	seenTiles := make(map[Position]bool)
	currentGuardState := board.guardState

	seenTiles[currentGuardState.pos] = true

	for {
		nextGuardState, ok := stepGuard(board, currentGuardState)
		if !ok {
			break
		}

		if currentGuardState.pos != nextGuardState.pos {
			seenTiles[nextGuardState.pos] = true
		}

		currentGuardState = nextGuardState
	}

	return seenTiles
}

func stepGuard(board Board, guardState GuardState) (GuardState, bool) {
	// step guard routine
	var nextTile Position
	switch guardState.facing {
	case 0:
		nextTile.x, nextTile.y = guardState.pos.x, guardState.pos.y-1
	case 1:
		nextTile.x, nextTile.y = guardState.pos.x+1, guardState.pos.y
	case 2:
		nextTile.x, nextTile.y = guardState.pos.x, guardState.pos.y+1
	case 3:
		nextTile.x, nextTile.y = guardState.pos.x-1, guardState.pos.y
	}

	if nextTile.x < 0 || nextTile.x >= board.width || nextTile.y < 0 || nextTile.y >= board.height {
		// next tile takes us out of bounds, exit
		return guardState, false
	}

	if _, ok := board.obstacles[nextTile]; ok {
		return GuardState{
			pos:    guardState.pos,
			facing: (guardState.facing + 1) % 4,
		}, true
	} else {
		return GuardState{
			pos:    nextTile,
			facing: guardState.facing,
		}, true
	}
}

func findLoops(board Board, guardPath map[Position]bool) int {
	// Hypothesis: we don't need to test all tiles in the map because we know the
	// guard's path already. we only need to test placing obstacles in the
	// original path.

	loopsFound := 0

	for pos := range guardPath {
		if pos == board.guardState.pos {
			continue
		}

		board.obstacles[pos] = true

		foundLoop := false
		newPath := make(map[GuardState]bool)
		currentGuardState := board.guardState
		newPath[currentGuardState] = true
		pathTrace := []GuardState{currentGuardState}
		for {
			nextGuardState, ok := stepGuard(board, currentGuardState)
			pathTrace = append(pathTrace, nextGuardState)
			if !ok {
				break
			}

			if _, ok := newPath[nextGuardState]; ok {
				foundLoop = true
				break
			} else {
				newPath[nextGuardState] = true
			}

			currentGuardState = nextGuardState

		}

		if foundLoop {
			loopsFound++
		}

		// reset obstacle
		delete(board.obstacles, pos)
	}

	return loopsFound
}
