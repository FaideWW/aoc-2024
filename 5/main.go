package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2024/lib"
)

func main() {
	input := io.ReadInputFile(os.Args[1])
	parts := io.TrimAndSplitBy(input, "\n\n")
	rules, updates := parts[0], parts[1]

	prereqs := computePrerequisitePages(io.TrimAndSplit(rules))

	correctResult, incorrectResult := computeUpdates(io.TrimAndSplit(updates), prereqs)
	fmt.Printf("valid updates sum: %d\n", correctResult)
	fmt.Printf("invalid updates sum: %d\n", incorrectResult)
}

func computePrerequisitePages(rules []string) map[int][]int {
	prereqs := make(map[int][]int)
	// for each rule:
	// the first page must be printed before the second page
	// prereqs is a map keyed by each page, and the value is a list of the pages that must be printed before it

	for _, ruleStr := range rules {
		pages := strings.Split(ruleStr, "|")
		requiredPage, _ := strconv.Atoi(pages[0])
		page, _ := strconv.Atoi(pages[1])

		if _, ok := prereqs[page]; !ok {
			prereqs[page] = []int{}
		}
		prereqs[page] = append(prereqs[page], requiredPage)
	}

	return prereqs
}

func computeUpdates(updates []string, pageTree map[int][]int) (int, int) {
	correctSum := 0
	incorrectSum := 0
	for _, updateStr := range updates {
		// for each page, we store the index of that page in a map.
		// after the map is built, we iterate over the pages again and
		// check that every prerequisite page (as dictated by pageTree)
		// does not have a higher index in the order map.
		pagesStr := strings.Split(updateStr, ",")
		pages := make([]int, len(pagesStr))
		orderMap := make(map[int]int)
		for i, pageStr := range pagesStr {
			pageNum, _ := strconv.Atoi(pageStr)
			pages[i] = pageNum
			orderMap[pageNum] = i
		}

		isCorrect := true

		for i, pageNum := range pages {
			prereqs, ok := pageTree[pageNum]
			if ok {
				for _, reqPage := range prereqs {
					if idx, ok := orderMap[reqPage]; ok && idx >= i {
						isCorrect = false
						break
					}
				}
			}
			if !isCorrect {
				break
			}
		}

		if isCorrect {
			correctSum += pages[len(pages)/2]
		} else {
			incorrectSum += fixInvalidUpdate(pages, orderMap, pageTree)
		}
	}

	return correctSum, incorrectSum
}

func fixInvalidUpdate(pages []int, orderMap map[int]int, pageTree map[int][]int) int {
	// for each pair of pages out of order, swap them.
	// repeat this until all the pages are in order
	noErrors := false
	for !noErrors {
		foundError := false
		for i, pageNum := range pages {
			prereqs, ok := pageTree[pageNum]
			if ok {
				for _, reqPage := range prereqs {
					if idx, ok := orderMap[reqPage]; ok && idx >= i {
						foundError = true
						pages[i], pages[idx] = pages[idx], pages[i]
						orderMap[pageNum] = idx
						orderMap[reqPage] = i
						break
					}
				}
			}
		}
		if !foundError {
			noErrors = true
		}
	}

	return pages[len(pages)/2]
}
