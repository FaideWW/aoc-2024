package main

import (
	"fmt"
	"os"
	"strconv"

	io "github.com/faideww/aoc-2024/lib"
)

type Position struct {
	x, y int
}

type Map struct {
	width, height int
	heights       map[Position]int
	trailheads    []Position
	trailCache    map[Position]int
}

func main() {
	input := io.ReadInputFile(os.Args[1])
	m := parseMap(input)

	result := computeTrailScores(m)
	fmt.Printf("trail score sum:%d\n", result)

	result2 := computeTrailRatings(m)
	fmt.Printf("trail rating sum:%d\n", result2)
}

func parseMap(input string) Map {
	lines := io.TrimAndSplit(input)
	width := len(lines[0])
	height := len(lines)
	heights := make(map[Position]int)
	trailheads := []Position{}
	for y, line := range lines {
		for x := range line {
			height, _ := strconv.Atoi(line[x : x+1])
			heights[Position{x, y}] = height
			if height == 0 {
				trailheads = append(trailheads, Position{x, y})
			}
		}
	}
	return Map{
		width:      width,
		height:     height,
		heights:    heights,
		trailheads: trailheads,
		trailCache: make(map[Position]int),
	}
}

func computeTrailScores(m Map) int {
	// Idea: given any position on the trail (not limited to just trailheads), we can compute the trail score for that position based on how many unique 9s are accessible. Then, the score for any trail leading into that position is that position's score + any other trails found. So we can solve this using dynamic programming: cache the trail scores for each position starting at each trailhead, and re-use the computed scores when we encounter those trails again in future searches.
	sum := 0
	for _, head := range m.trailheads {
		seenNines := make(map[Position]bool)
		sum += computeTrailScore(m, head, seenNines)

	}
	return sum
}

func computeTrailScore(m Map, p Position, seenNines map[Position]bool) int {
	currentHeight := m.heights[p]
	if currentHeight == 9 {
		if _, seen := seenNines[p]; !seen {
			seenNines[p] = true
			return 1
		} else {
			return 0
		}

	}
	if score, ok := m.trailCache[p]; ok {
		return score
	}

	neighbors := make([]Position, 0)
	if p.x > 0 {
		neighbors = append(neighbors, Position{p.x - 1, p.y})
	}
	if p.x < m.width-1 {
		neighbors = append(neighbors, Position{p.x + 1, p.y})
	}
	if p.y > 0 {
		neighbors = append(neighbors, Position{p.x, p.y - 1})
	}
	if p.y < m.height-1 {
		neighbors = append(neighbors, Position{p.x, p.y + 1})
	}

	score := 0
	for _, neighbor := range neighbors {
		if m.heights[neighbor]-currentHeight == 1 {
			score += computeTrailScore(m, neighbor, seenNines)
		}
	}

	return score
}

func computeTrailRatings(m Map) int {
	// similar to trail scores, but we care about all unique paths and not just unique endings.
	sum := 0
	for _, head := range m.trailheads {
		sum += computeTrailRating(m, head)
	}
	return sum
}

func computeTrailRating(m Map, p Position) int {
	currentHeight := m.heights[p]
	if currentHeight == 9 {
		return 1

	}
	if score, ok := m.trailCache[p]; ok {
		return score
	}

	neighbors := make([]Position, 0)
	if p.x > 0 {
		neighbors = append(neighbors, Position{p.x - 1, p.y})
	}
	if p.x < m.width-1 {
		neighbors = append(neighbors, Position{p.x + 1, p.y})
	}
	if p.y > 0 {
		neighbors = append(neighbors, Position{p.x, p.y - 1})
	}
	if p.y < m.height-1 {
		neighbors = append(neighbors, Position{p.x, p.y + 1})
	}

	score := 0
	for _, neighbor := range neighbors {
		if m.heights[neighbor]-currentHeight == 1 {
			score += computeTrailRating(m, neighbor)
		}
	}

	return score
}
