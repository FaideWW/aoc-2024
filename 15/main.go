package main

import (
	"fmt"
	"os"
	"strings"

	io "github.com/faideww/aoc-2024/lib"
)

type Position struct{ x, y int }

type Board struct {
	width, height int
	tiles         map[Position]rune
	robot         Position
	moves         string
	nextMove      int
}

func main() {
	input := io.ReadInputFile(os.Args[1])
	board := parseBoard(input)
	runRobot(board)
	score := scoreBoard(board)

	fmt.Printf("gps score: %d\n", score)

	board2 := wideParseBoard(input)
	wideRunRobot(board2)
	score2 := wideScoreBoard(board2)

	fmt.Printf("wide gps score: %d\n", score2)
}

func parseBoard(input string) Board {
	sections := io.TrimAndSplitBy(input, "\n\n")

	// parse map
	lines := io.TrimAndSplit(sections[0])
	width := len(lines[0])
	height := len(lines)
	tiles := make(map[Position]rune)

	var robotPos Position

	for y, line := range lines {
		for x, char := range line {
			pos := Position{x, y}
			tiles[pos] = char

			if char == '@' {
				robotPos = pos
			}
		}
	}

	moves := strings.Join(strings.Split(sections[1], "\n"), "")

	return Board{
		width:    width,
		height:   height,
		tiles:    tiles,
		robot:    robotPos,
		moves:    moves,
		nextMove: 0,
	}
}

func runRobot(board Board) {
	for board.nextMove < len(board.moves) {
		nextMove := board.moves[board.nextMove]
		board.nextMove++

		// attempt to move the robot
		var delta Position
		switch nextMove {
		case '^':
			delta.x = 0
			delta.y = -1
		case '>':
			delta.x = 1
			delta.y = 0
		case 'v':
			delta.x = 0
			delta.y = 1
		case '<':
			delta.x = -1
			delta.y = 0
		}

		current := board.robot
		dest := Position{current.x + delta.x, current.y + delta.y}
		canMove := true
		for board.tiles[dest] != '.' {
			// advance each tile one 'delta' forward from current
			if board.tiles[dest] == '#' {
				canMove = false
				break
			}
			current = dest
			dest = Position{current.x + delta.x, current.y + delta.y}
		}

		if canMove {
			// walk back from the dest to the robot, and shift all tiles forward

			current := dest
			for {
				// swap the tiles at current and current-reverseDelta
				prev := Position{current.x - delta.x, current.y - delta.y}
				board.tiles[current], board.tiles[prev] = board.tiles[prev], board.tiles[current]

				if board.tiles[current] == '@' {
					board.robot = current
					break
				}

				current = prev
			}
		}
	}
}

func printBoard(board Board) {
	for y := 0; y < board.height; y++ {
		for x := 0; x < board.width; x++ {
			fmt.Printf("%c", board.tiles[Position{x, y}])
		}
		fmt.Println()
	}
}

func scoreBoard(board Board) int {
	sum := 0
	for pos, tile := range board.tiles {
		if tile == 'O' {
			sum += pos.y*100 + pos.x
		}
	}

	return sum
}

func wideParseBoard(input string) Board {
	sections := io.TrimAndSplitBy(input, "\n\n")

	// parse map
	lines := io.TrimAndSplit(sections[0])
	width := len(lines[0]) * 2
	height := len(lines)
	tiles := make(map[Position]rune)

	var robotPos Position

	for y, line := range lines {
		for lx, char := range line {
			pos := Position{lx * 2, y}
			pos2 := Position{(lx * 2) + 1, y}
			switch char {
			case '#':
				tiles[pos] = '#'
				tiles[pos2] = '#'
			case '.':
				tiles[pos] = '.'
				tiles[pos2] = '.'
			case 'O':
				tiles[pos] = '['
				tiles[pos2] = ']'
			case '@':
				tiles[pos] = '@'
				tiles[pos2] = '.'
				robotPos = pos
			}
		}
	}

	moves := strings.Join(strings.Split(sections[1], "\n"), "")

	return Board{
		width:    width,
		height:   height,
		tiles:    tiles,
		robot:    robotPos,
		moves:    moves,
		nextMove: 0,
	}
}

func wideRunRobot(board Board) {
	for board.nextMove < len(board.moves) {
		nextMove := board.moves[board.nextMove]
		board.nextMove++

		// attempt to move the robot
		var delta Position
		switch nextMove {
		case '^':
			delta.x = 0
			delta.y = -1
		case '>':
			delta.x = 1
			delta.y = 0
		case 'v':
			delta.x = 0
			delta.y = 1
		case '<':
			delta.x = -1
			delta.y = 0
		}

		// walk forward, "staging" each block for movement. if any block is not movable, the entire move is cancelled
		if delta.y != 0 {
			// vertical movement is a special case; we need to make sure we deal with double-wide boxes appropriately
			if isTileMovable(board, board.robot, delta) {
				commitTileMovement(board, board.robot, delta)
				board.robot = Position{board.robot.x + delta.x, board.robot.y + delta.y}
			}
		} else {
			current := board.robot
			dest := Position{current.x + delta.x, current.y + delta.y}
			canMove := true
			for board.tiles[dest] != '.' {
				// advance each tile one 'delta' forward from current
				if board.tiles[dest] == '#' {
					canMove = false
					break
				}
				current = dest
				dest = Position{current.x + delta.x, current.y + delta.y}
			}

			if canMove {
				// walk back from the dest to the robot, and shift all tiles forward

				current := dest
				for {
					// swap the tiles at current and current-reverseDelta
					prev := Position{current.x - delta.x, current.y - delta.y}
					board.tiles[current], board.tiles[prev] = board.tiles[prev], board.tiles[current]

					if board.tiles[current] == '@' {
						board.robot = current
						break
					}

					current = prev
				}
			}
		}
		// printBoard(board)
	}
}

func isTileMovable(board Board, tile Position, delta Position) bool {
	nextPos := Position{tile.x + delta.x, tile.y + delta.y}
	if board.tiles[nextPos] == '.' {
		return true
	} else if board.tiles[nextPos] == '#' {
		return false
	} else if board.tiles[nextPos] == '[' {
		return isTileMovable(board, nextPos, delta) && isTileMovable(board, Position{nextPos.x + 1, nextPos.y}, delta)
	} else if board.tiles[nextPos] == ']' {
		return isTileMovable(board, Position{nextPos.x - 1, nextPos.y}, delta) && isTileMovable(board, nextPos, delta)
	}

	// should be unreachable
	return false
}

func commitTileMovement(board Board, tile Position, delta Position) {
	nextPos := Position{tile.x + delta.x, tile.y + delta.y}
	if board.tiles[nextPos] == '[' {
		commitTileMovement(board, nextPos, delta)
		commitTileMovement(board, Position{nextPos.x + 1, nextPos.y}, delta)
	} else if board.tiles[nextPos] == ']' {
		commitTileMovement(board, Position{nextPos.x - 1, nextPos.y}, delta)
		commitTileMovement(board, nextPos, delta)
	}
	board.tiles[tile], board.tiles[nextPos] = board.tiles[nextPos], board.tiles[tile]
}

func wideScoreBoard(board Board) int {
	sum := 0
	for pos, tile := range board.tiles {
		if tile == '[' {
			sum += pos.y*100 + pos.x
		}
	}

	return sum
}
