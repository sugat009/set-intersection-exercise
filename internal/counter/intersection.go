package counter

import (
	"errors"
)

// FindSetIntersection finds counts the intersection of keys between two key streams.
// Returns when both first and second channel is closed
func FindSetIntersection(first <-chan string, second <-chan string) (IntersectionResult, error) {
	if first == nil || second == nil {
		return IntersectionResult{}, errors.New("input channel cannot nil")
	}

	return findSetIntersection(first, second)
}

func findSetIntersection(first <-chan string, second <-chan string) (IntersectionResult, error) {
	firstDone, secondDone := make(chan int), make(chan int)
	var firstKeys, secondKeys map[string]int
	var firstTotalKeyCount, secondTotalKeyCount int

	go func() {
		firstKeys, firstTotalKeyCount = countKeys(first)
		close(firstDone)
	}()

	go func() {
		secondKeys, secondTotalKeyCount = countKeys(second)
		close(secondDone)
	}()

	<-firstDone
	<-secondDone

	distinctOverlap, totalOverlap := findOverlaps(firstKeys, secondKeys)

	result := IntersectionResult{
		First: FileResult{
			KeyCount:         firstTotalKeyCount,
			DistinctKeyCount: len(firstKeys),
		},
		Second: FileResult{
			KeyCount:         secondTotalKeyCount,
			DistinctKeyCount: len(secondKeys),
		},
		TotalOverlap:    totalOverlap,
		DistinctOverlap: distinctOverlap,
	}

	return result, nil
}

func countKeys(input <-chan string) (map[string]int, int) {
	res := make(map[string]int)
	totalCount := 0
	for {
		item, more := <-input
		if !more {
			break
		}

		if _, ok := res[item]; !ok {
			res[item] = 0
		}
		res[item]++
		totalCount++

	}
	return res, totalCount
}

func findOverlaps(firstKeys, secondKeys map[string]int) (int, int) {
	count := 0
	totalOverlap := 0

	for fk, fv := range firstKeys {
		if sv, ok := secondKeys[fk]; ok {
			count++

			if fv < sv {
				totalOverlap += fv
			} else {
				totalOverlap += sv
			}
		}
	}
	return count, totalOverlap
}

// IntersectionResult represents result of intersection count
type IntersectionResult struct {
	First           FileResult
	Second          FileResult
	TotalOverlap    int
	DistinctOverlap int
}

// FileResult represents result of a file key count
type FileResult struct {
	KeyCount         int
	DistinctKeyCount int
}
