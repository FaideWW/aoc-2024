package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"

	io "github.com/faideww/aoc-2024/lib"
)

const PART2_PREFIX = 10000000000000

type Position struct{ x, y int }

type Game struct {
	prize   Position
	aButton Position
	bButton Position
}

func main() {
	input := io.ReadInputFile(os.Args[1])

	games := parseGames(input)

	sum := 0
	for _, g := range games {
		cost := findMinTokenCountLarge(g)
		if cost > 0 {
			sum += cost
		}
	}
	fmt.Printf("total cost: %d\n", sum)

	sum = 0
	for _, g := range games {
		g.prize.x += PART2_PREFIX
		g.prize.y += PART2_PREFIX
		cost := findMinTokenCountLarge(g)
		if cost > 0 {
			sum += cost
		}
	}

	fmt.Printf("total cost (part 2): %d\n", sum)
}

func parseGames(input string) []Game {
	aButtonRegex := regexp.MustCompile("Button A: X+(.+), Y+(.+)$")
	bButtonRegex := regexp.MustCompile("Button B: X+(.+), Y+(.+)$")
	prizeRegex := regexp.MustCompile("Prize: X=(.+), Y=(.+)$")
	gameStrs := io.TrimAndSplitBy(input, "\n\n")
	games := []Game{}
	for _, gameStr := range gameStrs {
		lines := io.TrimAndSplit(gameStr)
		aButton := aButtonRegex.FindStringSubmatch(lines[0])
		bButton := bButtonRegex.FindStringSubmatch(lines[1])
		prize := prizeRegex.FindStringSubmatch(lines[2])
		aX, _ := strconv.Atoi(aButton[1])
		aY, _ := strconv.Atoi(aButton[2])
		bX, _ := strconv.Atoi(bButton[1])
		bY, _ := strconv.Atoi(bButton[2])
		prizeX, _ := strconv.Atoi(prize[1])
		prizeY, _ := strconv.Atoi(prize[2])

		games = append(games, Game{
			prize:   Position{prizeX, prizeY},
			aButton: Position{aX, aY},
			bButton: Position{bX, bY},
		})
	}
	return games
}

// Using uniform-cost search: we need to pathfind from 0,0 to the
// target position using button presses that advance us forward.
// Additionally, each edge is weighted based on which button is
// being pressed (A presses cost 3, B presses cost 1). Once we
// have a path, we simply count the number of presses for each
// button and return the result.
func findMinTokenCount(game Game) int {
	origin := Position{0, 0}
	pq := io.PriorityQueue[Position]{}
	heap.Push(&pq, &io.PQItem[Position]{
		Value:    origin,
		Priority: 0,
	})

	expanded := make(map[Position]bool)

	for pq.Len() > 0 {
		currentNode := heap.Pop(&pq).(*io.PQItem[Position])
		current := currentNode.Value
		// fmt.Printf("current: %v\n", current)
		if current == game.prize {
			return currentNode.Priority // the priority will be the total cost
		}
		expanded[current] = true

		if current.x > game.prize.x || current.y > game.prize.y {
			// buttons strictly only advance position forward, never backward. so if we've gone past the goal we can stop.
			continue
		}

		// push A button
		currentPlusA := Position{current.x + game.aButton.x, current.y + game.aButton.y}
		_, inExpanded := expanded[currentPlusA]
		if !inExpanded {
			heap.Push(&pq, &io.PQItem[Position]{
				Value:    currentPlusA,
				Priority: currentNode.Priority + 3, // A button costs 3
			})
		}

		// push B button
		currentPlusB := Position{current.x + game.bButton.x, current.y + game.bButton.y}
		_, inExpanded = expanded[currentPlusB]
		if !inExpanded {
			heap.Push(&pq, &io.PQItem[Position]{
				Value:    currentPlusB,
				Priority: currentNode.Priority + 1, // B button costs 1
			})
		}
	}

	return -1
}

// Part 2: Find the much larger targets with gaussian elimination.
//
// Since we are only ever working with a 2x3 matrix we can hand-solve it
// instead of having to implement a full gaussian elimination solver.
func findMinTokenCountLarge(game Game) int {
	button1, button2 := game.aButton, game.bButton
	swapped := false
	if button1.x < button2.x {
		// swap the buttons so that we end up with a positive solution
		swapped = true
		button1, button2 = button2, button1
	}

	aFactor := float64(button1.y) / float64(button1.x)

	b_y := float64(button2.y) - (float64(button2.x) * aFactor)
	p_y := float64(game.prize.y) - (float64(game.prize.x) * aFactor)

	j := p_y / b_y
	i := (float64(game.prize.x) - (float64(button2.x) * j)) / float64(button1.x)

	// check that both i and j resolve to a whole number (if not, there is no solution)
	if math.Abs(j-math.Round(j)) > 0.001 || math.Abs(i-math.Round(i)) > 0.001 {
		return -1
	}

	jInt := int(math.Round(j))
	iInt := int(math.Round(i))

	tokenCost := 0
	// un-swap the results at the end
	if !swapped {
		tokenCost = iInt*3 + jInt
	} else {
		tokenCost = jInt*3 + iInt
	}

	return tokenCost
}
