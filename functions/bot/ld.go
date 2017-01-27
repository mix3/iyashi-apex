package main

import "fmt"

type LevenshteinDistance map[string]int

func (ld LevenshteinDistance) get(x, y int) int {
	return ld[fmt.Sprintf("%d:%d", x, y)]
}

func (ld LevenshteinDistance) set(x, y, v int) {
	ld[fmt.Sprintf("%d:%d", x, y)] = v
}

func (ld LevenshteinDistance) min(first int, more ...int) int {
	min := first
	for _, v := range more {
		if v < min {
			min = v
		}
	}
	return min
}

func (ld LevenshteinDistance) diff(a, b rune) int {
	diff := 0
	if a != b {
		diff = 1
	}
	return diff
}

func (ld LevenshteinDistance) Distance(s1, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)

	for i := 0; i <= len(r1); i++ {
		ld.set(i, 0, i)
	}
	for j := 0; j <= len(r1); j++ {
		ld.set(0, j, j)
	}

	for i, v1 := range r1 {
		for j, v2 := range r2 {
			ld.set(i+1, j+1, ld.min(
				ld.get(i, j)+ld.diff(v1, v2),
				ld.get(i, j+1)+1,
				ld.get(i+1, j)+1,
			))
		}
	}

	return ld.get(len(r1), len(r2))
}
