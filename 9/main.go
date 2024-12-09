package main

import (
	"fmt"
	"os"
	"strconv"

	io "github.com/faideww/aoc-2024/lib"
)

type File struct {
	id    int
	start int
	size  int
}

type Disk struct {
	length int
	files  []File
}

func main() {
	input := io.ReadInputFile(os.Args[1])

	disk := parseDisk(input)
	checksum := computeCompactedChecksum(disk)

	fmt.Printf("checksum: %d\n", checksum)

	defragmentDisk(&disk)
	defraggedChecksum := computeChecksum(disk)
	fmt.Printf("defragged checksum: %d\n", defraggedChecksum)
}

func parseDisk(input string) Disk {
	files := []File{}

	currentOffset := 0
	currentId := 0
	fileNext := true
	for i := 0; i < len(input); i++ {
		value, _ := strconv.Atoi(input[i : i+1])
		if !fileNext {
			fileNext = true
		} else {
			files = append(files, File{
				id:    currentId,
				start: currentOffset,
				size:  value,
			})
			currentId++
			fileNext = false
		}
		currentOffset += value

	}
	return Disk{
		length: currentOffset,
		files:  files,
	}
}

func printDisk(disk Disk) {
	lastOffset := 0
	for _, f := range disk.files {
		for ; lastOffset < f.start; lastOffset++ {
			fmt.Printf(".")
		}
		for ; lastOffset < f.start+f.size; lastOffset++ {
			fmt.Printf("%d", f.id)
		}
	}

	fmt.Printf("\n")
}

func computeCompactedChecksum(disk Disk) int {
	// key idea: keep a reverse cursor that starts at the
	// end of the last file, and a forward cursor that
	// starts at the front. while computing the checksum,
	// move the forward cursor forward until we read an
	// empty space, then read from the reverse cursor (and
	// move it backward accordingly, skipping spaces)

	lastFileIdx := len(disk.files) - 1
	lastFile := disk.files[lastFileIdx]
	reverseCursor := lastFile.start + lastFile.size - 1

	checksum := 0

	currentFileIdx := 0
	currentFile := disk.files[currentFileIdx]

	for forwardCursor := 0; forwardCursor <= reverseCursor; forwardCursor++ {
		// fmt.Printf("forwardCursor:%d currentFile: %+v\n", forwardCursor, currentFile)
		if forwardCursor < currentFile.start {
			// read from the reverse cursor
			value := lastFile.id
			checksum += forwardCursor * value
			reverseCursor--
			if reverseCursor < lastFile.start {
				// if we reach the beginning of the file, move to the end of the previous one
				lastFileIdx--
				lastFile = disk.files[lastFileIdx]
				reverseCursor = lastFile.start + lastFile.size - 1
			}

		} else {
			// read from the forward cursor
			value := currentFile.id
			checksum += forwardCursor * value

			if forwardCursor == currentFile.start+currentFile.size-1 {
				currentFileIdx++
				currentFile = disk.files[currentFileIdx]
			}
		}
	}

	fmt.Printf("\n")

	return checksum
}

// Part 2

func defragmentDisk(disk *Disk) {
	moved := make(map[int]bool)
	for i := len(disk.files) - 1; i >= 0; i-- {
		targetFile := disk.files[i]
		if _, ok := moved[targetFile.id]; ok {
			// we've already tried to move this file, don't try again
			continue
		}

		// find a space that fits the file (if exists)
		for nextFile := 0; nextFile < len(disk.files); nextFile++ {
			var gapStart, gapSize int
			if nextFile == 0 {
				gapStart = 0
				gapSize = disk.files[nextFile].start
			} else {
				gapStart = disk.files[nextFile-1].start + disk.files[nextFile-1].size
				gapSize = disk.files[nextFile].start - gapStart
			}

			if gapSize >= targetFile.size && gapStart < targetFile.start {
				moved[targetFile.id] = true
				targetFile.start = gapStart
				// re-order the disk (hacky, but i don't want to write a binary tree implementation lol)

				for j := i - 1; j >= nextFile; j-- {
					disk.files[j+1] = disk.files[j]
				}
				disk.files[nextFile] = targetFile
				i++
				break
			}
		}
	}
}

func computeChecksum(disk Disk) (checksum int) {
	currentFileIdx := 0
	for cursor := 0; cursor < disk.length; cursor++ {
		if currentFileIdx >= len(disk.files) {
			continue
		}
		currentFile := disk.files[currentFileIdx]
		if cursor < currentFile.start {
			continue
		}
		value := currentFile.id
		checksum += cursor * value

		if cursor == currentFile.start+currentFile.size-1 {
			currentFileIdx++
		}
	}
	return
}
