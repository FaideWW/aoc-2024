package io

import (
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return strings.TrimSpace(string(dat))
}

func TrimAndSplit(input string) []string {
	return strings.Split(strings.ReplaceAll(strings.TrimSpace(input), "\r\n", "\n"), "\n")
}

func TrimAndSplitBy(input string, delimiter string) []string {
	return strings.Split(strings.ReplaceAll(strings.TrimSpace(input), "\r\n", "\n"), delimiter)
}

// PriorityQueue implementation from https://pkg.go.dev/container/heap#example__priorityQueue

type PQItem[T any] struct {
	Value    T
	Priority int
	Index    int
}

type PriorityQueue[T any] []*PQItem[T]

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue[T]) Push(x any) {
	n := len(*pq)
	item := x.(*PQItem[T])
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// identical to the above but sorts in ascending order
type PriorityQueueAsc[T any] []*PQItem[T]

func (pq PriorityQueueAsc[T]) Len() int { return len(pq) }

func (pq PriorityQueueAsc[T]) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueueAsc[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueueAsc[T]) Push(x any) {
	n := len(*pq)
	item := x.(*PQItem[T])
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueueAsc[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func PowInt(x, y int) int {
	if y == 0 {
		return 1
	}
	if y == 1 {
		return x
	}
	result := x
	for i := 1; i < y; i++ {
		result *= x
	}

	return result
}
