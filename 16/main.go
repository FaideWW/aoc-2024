package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"

	io "github.com/faideww/aoc-2024/lib"
)

type Position struct {
	x, y int
}

const (
	NORTH = iota
	EAST
	SOUTH
	WEST
)

type Position2 struct {
	x, y, dir int
}

type Board struct {
	width, height int
	tiles         map[Position]rune
	start         Position
	goal          Position
}

func main() {
	input := io.ReadInputFile(os.Args[1])
	board := parseBoard(input)

	score := findCheapestRoute(board)
	fmt.Printf("score: %d\n", score)

	tilesOnPath := findAllCheapestRoutes(board)
	fmt.Printf("tilesOnPath: %d\n", tilesOnPath)
}

func parseBoard(input string) Board {
	lines := io.TrimAndSplit(input)
	width := len(lines[0])
	height := len(lines)
	tiles := make(map[Position]rune)
	var start, goal Position
	for y, line := range lines {
		for x, char := range line {
			pos := Position{x, y}
			if char == 'S' {
				start = pos
			} else if char == 'E' {
				goal = pos
			}

			tiles[pos] = char
		}
	}
	return Board{
		width,
		height,
		tiles,
		start,
		goal,
	}
}

func findCheapestRoute(board Board) int {
	frontier := io.PriorityQueueAsc[Position]{}
	heap.Push(&frontier, &io.PQItem[Position]{
		Value:    board.start,
		Priority: 0,
	})

	cameFrom := make(map[Position]Position)

	// Reindeer start facing east, so we say they "came from" the tile to the left
	cameFrom[board.start] = Position{board.start.x - 1, board.start.y}

	costSoFar := make(map[Position]int)
	costSoFar[board.start] = 0

	for frontier.Len() > 0 {
		currentNode := heap.Pop(&frontier).(*io.PQItem[Position])
		current := currentNode.Value

		lastDirection := Position{current.x - cameFrom[current].x, current.y - cameFrom[current].y}

		neighbors := findNeighbors(board.tiles, current, lastDirection)
		for _, neighbor := range neighbors {
			newCost := costSoFar[current] + neighbor.cost
			if oldCost, ok := costSoFar[neighbor.pos]; !ok || newCost < oldCost {
				costSoFar[neighbor.pos] = newCost
				cameFrom[neighbor.pos] = current
				heap.Push(&frontier, &io.PQItem[Position]{
					Value:    neighbor.pos,
					Priority: newCost,
				})
			}
		}
	}

	return costSoFar[board.goal]
}

type PosCost struct {
	pos  Position
	cost int
}

func findNeighbors(tiles map[Position]rune, pos Position, lastDirection Position) []PosCost {
	neighbors := make([]PosCost, 0)

	dirs := []Position{
		Position{0, -1},
		Position{1, 0},
		Position{0, 1},
		Position{-1, 0},
	}

	for _, dir := range dirs {
		nextPos := Position{pos.x + dir.x, pos.y + dir.y}
		negativeDir := Position{-dir.x, -dir.y}
		if negativeDir == lastDirection {
			// we don't need to check 180-degree turns
			continue
		}
		if tiles[nextPos] != '#' {
			if dir == lastDirection && tiles[nextPos] != '#' {
				neighbors = append(neighbors, PosCost{nextPos, 1})
			} else {
				neighbors = append(neighbors, PosCost{nextPos, 1001})
			}
		}
	}

	return neighbors
}

type RouteNode struct {
	pos  Position2
	cost int
	path []Position2
}

func findAllCheapestRoutes(board Board) int {
	frontier := io.PriorityQueueAsc[RouteNode]{}
	start := Position2{board.start.x, board.start.y, EAST}
	heap.Push(&frontier, &io.PQItem[RouteNode]{
		Value:    RouteNode{start, 0, []Position2{start}},
		Priority: 0,
	})

	visited := make(map[Position2]int)
	onPathTiles := make(map[Position]bool)

	maxCost := math.MaxInt

	for frontier.Len() > 0 {
		currentNode := heap.Pop(&frontier).(*io.PQItem[RouteNode])
		current := currentNode.Value

		if current.cost > maxCost {
			break
		}

		if lastCost, ok := visited[current.pos]; ok && lastCost < current.cost {
			continue
		}

		visited[current.pos] = current.cost

		if current.pos.x == board.goal.x && current.pos.y == board.goal.y {
			maxCost = current.cost
			for _, tile := range current.path {
				onPathTiles[Position{tile.x, tile.y}] = true
			}
		}

		neighbors := findNeighbors2(board.tiles, current.pos)
		for _, neighbor := range neighbors {
			newCost := current.cost + neighbor.cost
			newPath := make([]Position2, len(current.path)+1)
			copy(newPath, current.path)
			newPath[len(newPath)-1] = neighbor.pos
			heap.Push(&frontier, &io.PQItem[RouteNode]{
				Value:    RouteNode{neighbor.pos, newCost, newPath},
				Priority: newCost,
			})
		}
	}

	// printBoard(board, onPathTiles)

	return len(onPathTiles)
}

type PosCost2 struct {
	pos  Position2
	cost int
}

func findNeighbors2(tiles map[Position]rune, pos Position2) []PosCost2 {
	var dx, dy int
	switch pos.dir {
	case NORTH:
		dx = 0
		dy = -1
	case EAST:
		dx = 1
		dy = 0
	case SOUTH:
		dx = 0
		dy = 1
	case WEST:
		dx = -1
		dy = 0
	}

	neighbors := make([]PosCost2, 0)

	nextPos := Position{pos.x + dx, pos.y + dy}
	if tiles[nextPos] != '#' {
		neighbors = append(neighbors, PosCost2{Position2{nextPos.x, nextPos.y, pos.dir}, 1})
	}

	dirLeft := (((pos.dir - 1) % 4) + 4) % 4
	dirRight := (pos.dir + 1) % 4

	neighbors = append(neighbors, PosCost2{Position2{pos.x, pos.y, dirLeft}, 1000})
	neighbors = append(neighbors, PosCost2{Position2{pos.x, pos.y, dirRight}, 1000})

	return neighbors
}

func printBoard(board Board, onPathTiles map[Position]bool) {
	for y := 0; y < board.height; y++ {
		for x := 0; x < board.width; x++ {
			pos := Position{x, y}
			if _, ok := onPathTiles[pos]; ok {
				fmt.Printf("O")
			} else {
				fmt.Printf("%c", board.tiles[pos])
			}
		}
		fmt.Println()
	}
}
