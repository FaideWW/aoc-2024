package main

import (
	"fmt"
	"os"

	io "github.com/faideww/aoc-2024/lib"
)

type Position struct {
	x, y int
}

type Board struct {
	width  int
	height int
	plants map[Position]rune
}

type Region struct {
	plant     rune
	area      int
	perimeter int
	cells     map[Position]bool
}

func main() {
	input := io.ReadInputFile(os.Args[1])

	board := parseBoard(input)

	regions := findRegions(board)

	result := calcTotalCost(regions)
	fmt.Printf("cost: %d\n", result)

	sum := 0
	for _, region := range regions {
		corners := countSides(board, region)
		sum += corners * region.area
	}

	fmt.Printf("bulk cost:%d\n", sum)
}

func parseBoard(input string) Board {
	lines := io.TrimAndSplit(input)
	height := len(lines)
	width := len(lines[0])

	plants := make(map[Position]rune)

	for y, line := range lines {
		for x, char := range line {
			plants[Position{x, y}] = char
		}
	}
	return Board{width, height, plants}
}

func findRegions(board Board) []Region {
	newRegionFrontier := make([]Position, 0)
	newRegionFrontier = append(newRegionFrontier, Position{0, 0})

	// seenPlants := make(map[Position]bool)
	seenRegions := make(map[Position]bool)

	regions := make([]Region, 0)

	// starting from 0,0 crawl the board adding neighbors to one of two lists:
	// - if it's the same plant type, add it to the sameRegionFrontier
	// - if it's a different plant type, add it to the newRegionFrontier

	for len(newRegionFrontier) > 0 {
		regionStart := newRegionFrontier[0]
		currentPlant := board.plants[regionStart]
		newRegionFrontier = newRegionFrontier[1:]

		if _, seen := seenRegions[regionStart]; seen {
			// if we've already encountered this plant as part of a region search, skip it
			continue
		}

		area := 0
		perimeter := 0

		sameRegionFrontier := make([]Position, 0)
		sameRegionFrontier = append(sameRegionFrontier, regionStart)

		seenInRegion := make(map[Position]bool)

		for len(sameRegionFrontier) > 0 {
			currentPos := sameRegionFrontier[0]
			sameRegionFrontier = sameRegionFrontier[1:]
			if _, seen := seenInRegion[currentPos]; seen {
				continue
			}
			seenInRegion[currentPos] = true
			seenRegions[currentPos] = true

			// find neighboring plants
			neighbors := findNeighbors(board, currentPos)
			neighboringSamePlants := 0

			for _, neighbor := range neighbors {
				if board.plants[neighbor] == currentPlant {
					neighboringSamePlants++
				}
				if board.plants[neighbor] == currentPlant {
					sameRegionFrontier = append(sameRegionFrontier, neighbor)
				} else {
					newRegionFrontier = append(newRegionFrontier, neighbor)
				}
			}

			area++
			perimeter += 4 - neighboringSamePlants
		}

		regions = append(regions, Region{
			plant:     currentPlant,
			area:      area,
			perimeter: perimeter,
			cells:     seenInRegion,
		})

	}

	return regions
}

func findNeighbors(board Board, pos Position) []Position {
	neighbors := make([]Position, 0)

	if pos.x > 0 {
		neighbors = append(neighbors, Position{pos.x - 1, pos.y})
	}
	if pos.x < board.width-1 {
		neighbors = append(neighbors, Position{pos.x + 1, pos.y})
	}
	if pos.y > 0 {
		neighbors = append(neighbors, Position{pos.x, pos.y - 1})
	}
	if pos.y < board.height-1 {
		neighbors = append(neighbors, Position{pos.x, pos.y + 1})
	}

	return neighbors
}

func calcTotalCost(regions []Region) int {
	sum := 0

	for _, region := range regions {
		sum += region.area * region.perimeter
	}

	return sum
}

func countSides(board Board, region Region) int {
	// the number of sides of a region is equal to the number of corners in that
	// region. region corners can be found by iterating over all cell corners in
	// the region and identifying one of the following cases:
	// 1. a cell corner is a convex region corner if it has 3 non-region cells
	//    adjacent to it
	// 2. a cell corner is a concave region corner if it has exactly 1 non-region
	//    cell adjacent to it
	// note that inner corners will get counted multiple times (3 to be exact).
	// we could fix this, but a simpler hack is to divide the inner corner count
	// by 3 at the end.

	corners := 0

	innerCorners := 0

	for cell := range region.cells {
		topLeftNeighbors := 0
		topRightNeighbors := 0
		bottomLeftNeighbors := 0
		bottomRightNeighbors := 0
		var tr, tl, br, bl, top, left, right, bottom bool
		if plant, ok := board.plants[Position{cell.x - 1, cell.y - 1}]; !ok || plant != region.plant {
			tl = true
			topLeftNeighbors++
		}
		if plant, ok := board.plants[Position{cell.x, cell.y - 1}]; !ok || plant != region.plant {
			top = true
			topLeftNeighbors++
			topRightNeighbors++
		}
		if plant, ok := board.plants[Position{cell.x + 1, cell.y - 1}]; !ok || plant != region.plant {
			tr = true
			topRightNeighbors++
		}

		if plant, ok := board.plants[Position{cell.x + 1, cell.y}]; !ok || plant != region.plant {
			right = true
			topRightNeighbors++
			bottomRightNeighbors++
		}
		if plant, ok := board.plants[Position{cell.x + 1, cell.y + 1}]; !ok || plant != region.plant {
			br = true
			bottomRightNeighbors++
		}
		if plant, ok := board.plants[Position{cell.x, cell.y + 1}]; !ok || plant != region.plant {
			bottom = true
			bottomRightNeighbors++
			bottomLeftNeighbors++
		}
		if plant, ok := board.plants[Position{cell.x - 1, cell.y + 1}]; !ok || plant != region.plant {
			bl = true
			bottomLeftNeighbors++
		}
		if plant, ok := board.plants[Position{cell.x - 1, cell.y}]; !ok || plant != region.plant {
			left = true
			bottomLeftNeighbors++
			topLeftNeighbors++
		}

		if topLeftNeighbors == 3 {
			corners++
		} else if topLeftNeighbors == 1 {
			innerCorners++
		}
		if topRightNeighbors == 3 {
			corners++
		} else if topRightNeighbors == 1 {
			innerCorners++
		}
		if bottomLeftNeighbors == 3 {
			corners++
		} else if bottomLeftNeighbors == 1 {
			innerCorners++
		}
		if bottomRightNeighbors == 3 {
			corners++
		} else if bottomRightNeighbors == 1 {
			innerCorners++
		}

		// special case: if two inner corners intersect in a checkerboard pattern,
		// we need to detect this separately. these don't need to be deduplicated,
		// so we add them directly to corners
		if top && !tr && right {
			corners++
		}
		if right && !br && bottom {
			corners++
		}
		if bottom && !bl && left {
			corners++
		}
		if left && !tl && top {
			corners++
		}

	}

	return corners + (innerCorners / 3)
}
