package main

import (
	// "bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	io "github.com/faideww/aoc-2024/lib"
)

const BOARD_WIDTH = 101
const BOARD_HEIGHT = 103

type Vec2 struct{ x, y int }

type Robot struct {
	position Vec2
	velocity Vec2
}

type Board struct {
	width, height int
	robots        []Robot
}

func main() {
	input := io.ReadInputFile(os.Args[1])
	board := parseBoard(input)

	advanceBoard(&board, 100)

	result := computeSafetyFactor(board)
	fmt.Printf("safety factor: %d\n", result)

	// part 2
	board = parseBoard(input)
	i := 1
	for {
		advanceBoard(&board, 1)
		if isBoardUnique(board) {
			// This happens to work for my input, but if there are other
			// configurations of robots in all-unique positions, we can simply
			// continue looping until we find it.
			fmt.Printf("turn %d\n", i)
			printBoard(board)
			break
		}
		i++
	}
}

func parseBoard(input string) Board {
	botRegex := regexp.MustCompile("p=(.+),(.+) v=(.+),(.+)")
	lines := io.TrimAndSplit(input)

	robots := make([]Robot, 0)

	for _, line := range lines {
		match := botRegex.FindStringSubmatch(line)
		px, _ := strconv.Atoi(match[1])
		py, _ := strconv.Atoi(match[2])
		vx, _ := strconv.Atoi(match[3])
		vy, _ := strconv.Atoi(match[4])
		r := Robot{
			Vec2{px, py},
			Vec2{vx, vy},
		}

		robots = append(robots, r)

	}
	return Board{
		width:  BOARD_WIDTH,
		height: BOARD_HEIGHT,
		robots: robots,
	}
}

func advanceBoard(board *Board, steps int) {
	for rIdx, r := range board.robots {
		deltaX := r.velocity.x * steps
		deltaY := r.velocity.y * steps
		nextX := (((r.position.x + deltaX) % board.width) + board.width) % board.width
		nextY := (((r.position.y + deltaY) % board.height) + board.height) % board.height
		board.robots[rIdx].position.x = nextX
		board.robots[rIdx].position.y = nextY
	}
}

func computeSafetyFactor(board Board) int {
	tl := 0
	tr := 0
	bl := 0
	br := 0
	for _, r := range board.robots {
		if r.position.x < board.width/2 && r.position.y < board.height/2 {
			tl++
		} else if r.position.x > board.width/2 && r.position.y < board.height/2 {
			tr++
		} else if r.position.x < board.width/2 && r.position.y > board.height/2 {
			bl++
		} else if r.position.x > board.width/2 && r.position.y > board.height/2 {
			br++
		}
	}
	return tl * tr * bl * br
}

func isBoardUnique(board Board) bool {
	robots := make(map[int]map[int]int)
	for _, r := range board.robots {
		if _, ok := robots[r.position.y]; !ok {
			robots[r.position.y] = make(map[int]int)
		}

		if robots[r.position.y][r.position.x] > 0 {
			return false
		}

		robots[r.position.y][r.position.x]++
	}

	return true
}

func printBoard(board Board) {
	robots := make(map[int]map[int]int)

	for _, r := range board.robots {
		if _, ok := robots[r.position.y]; !ok {
			robots[r.position.y] = make(map[int]int)
		}

		robots[r.position.y][r.position.x]++
	}

	for y := 0; y < board.height; y++ {
		for x := 0; x < board.width; x++ {
			if count, ok := robots[y][x]; ok && count > 0 {
				fmt.Printf("%d", count)
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}
