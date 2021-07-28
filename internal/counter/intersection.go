package counter

import (
	"errors"
	"sync"

	"github.com/rickyshrestha/infosum-interview-exercise/internal/pool"
)

// FindSetIntersection finds counts the intersection of keys between two key streams.
// Returns when both first and second channel is closed
func FindSetIntersection(first <-chan string, second <-chan string) (IntersectionResult, error) {
	if first == nil || second == nil {
		return IntersectionResult{}, errors.New("input channel cannot nil")
	}

	// find out if any channels are closed
	return findSetIntersection(first, second)
}

func findSetIntersection(first <-chan string, second <-chan string) (IntersectionResult, error) {
	var firstKeys, secondKeys map[string]int
	var firstTotalKeyCount, secondTotalKeyCount int

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		firstKeys, firstTotalKeyCount = countKeys(first)
		wg.Done()
	}()

	go func() {
		secondKeys, secondTotalKeyCount = countKeys(second)
		wg.Done()
	}()

	wg.Wait()

	//distinctOverlap, totalOverlap := findOverlapsUsingWorkerPool(firstKeys, secondKeys, 1024)
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

	var noMore bool
	for {
		select {
		case item, more := <-input:
			if !more {
				noMore = true
				break
			}

			if _, ok := res[item]; !ok {
				res[item] = 0
			}
			res[item]++
			totalCount++
		default:
		}

		if noMore {
			break
		}
	}
	return res, totalCount
}

func findOverlaps(firstKeys, secondKeys map[string]int) (int, int) {
	distinctOverlaps := 0
	totalOverlaps := 0

	for fk, fv := range firstKeys {
		if sv, ok := secondKeys[fk]; ok {
			distinctOverlaps++

			if fv < sv {
				totalOverlaps += fv
			} else {
				totalOverlaps += sv
			}
		}
	}
	return distinctOverlaps, totalOverlaps
}

func getOverlapCalcFunc(fk string, fv int, secondKeys map[string]int, distinctCounter, totalCounter *int) func() {
	return func() {
		if sv, ok := secondKeys[fk]; ok {
			*distinctCounter++

			if fv < sv {
				*totalCounter += fv
			} else {
				*totalCounter += sv
			}
		}
	}
}

func findOverlapsUsingWorkerPool(firstKeys, secondKeys map[string]int, workerPoolSize int) (int, int) {
	count := 0
	totalOverlap := 0

	workerPool := pool.NewWorkerPool(workerPoolSize)

	for fk, fv := range firstKeys {
		workerPool.AddWorker(getOverlapCalcFunc(fk, fv, secondKeys, &count, &totalOverlap))
	}

	workerPool.Wait()

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
