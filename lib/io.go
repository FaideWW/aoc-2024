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
