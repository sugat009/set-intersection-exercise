package counter

import "sync"

// FindSetIntersection finds counts the intersection of keys between two key streams.
// Returns when both first and second channel is closed
func FindSetIntersection(first <-chan string, second <-chan string) (*IntersectionResult, error) {
	return findSetIntersectionV3(first, second)
}

func findSetIntersectionV1(first <-chan string, second <-chan string) (*IntersectionResult, error) {
	firstKeys := make(map[string]int)
	secondKeys := make(map[string]int)

	firstTotalKeyCount, secondTotalKeyCount := 0, 0

	var firstDone, secondDone bool
	for !firstDone || !secondDone {
		select {
		case firstVal, more := <-first:
			if more {
				if _, ok := firstKeys[firstVal]; !ok {
					firstKeys[firstVal] = 0
				}
				firstKeys[firstVal]++
				firstTotalKeyCount++
			} else {
				firstDone = true
			}
		case secondVal, more := <-second:
			if more {
				if _, ok := secondKeys[secondVal]; !ok {
					secondKeys[secondVal] = 0
				}
				secondKeys[secondVal]++
				secondTotalKeyCount++
			} else {
				secondDone = true
			}
		}
	}

	distinctOverlap, totalOverlap := findOverlaps(firstKeys, secondKeys)

	result := &IntersectionResult{
		First: &FileResult{
			KeyCount:         firstTotalKeyCount,
			DistinctKeyCount: len(firstKeys),
		},
		Second: &FileResult{
			KeyCount:         secondTotalKeyCount,
			DistinctKeyCount: len(secondKeys),
		},
		TotalOverlap:    totalOverlap,
		DistinctOverlap: distinctOverlap,
	}

	return result, nil
}

func findSetIntersectionV2(first <-chan string, second <-chan string) (*IntersectionResult, error) {
	firstKeys, firstTotalKeyCount := countKeys(first)
	secondKeys, secondTotalKeyCount := countKeys(second)

	distinctOverlap, totalOverlap := findOverlaps(firstKeys, secondKeys)

	result := &IntersectionResult{
		First: &FileResult{
			KeyCount:         firstTotalKeyCount,
			DistinctKeyCount: len(firstKeys),
		},
		Second: &FileResult{
			KeyCount:         secondTotalKeyCount,
			DistinctKeyCount: len(secondKeys),
		},
		TotalOverlap:    totalOverlap,
		DistinctOverlap: distinctOverlap,
	}

	return result, nil
}

func findSetIntersectionV3(first <-chan string, second <-chan string) (*IntersectionResult, error) {
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

	result := &IntersectionResult{
		First: &FileResult{
			KeyCount:         firstTotalKeyCount,
			DistinctKeyCount: len(firstKeys),
		},
		Second: &FileResult{
			KeyCount:         secondTotalKeyCount,
			DistinctKeyCount: len(secondKeys),
		},
		TotalOverlap:    totalOverlap,
		DistinctOverlap: distinctOverlap,
	}

	return result, nil
}

// blocks until both channel is closed
func findSetIntersectionV4(first <-chan string, second <-chan string) (*IntersectionResult, error) {
	var firstTotalKeys, firstDistinctKeys, secondTotalKeys, secondDistinctKeys, totalOverlap, distinctOverlap int

	firstMap := make(map[string]int)
	secondMap := make(map[string]int)

	var firstDone, secondDone bool
	for !firstDone || !secondDone {
		select {
		case firstVal, more := <-first:
			if more {
				if _, ok := firstMap[firstVal]; !ok {
					firstMap[firstVal] = 0
					firstDistinctKeys++

					if _, ok := secondMap[firstVal]; ok {
						distinctOverlap++
					}
				}

				if countOnSecond, ok := secondMap[firstVal]; ok {
					previousVal := findMinimum(firstMap[firstVal], countOnSecond)
					newVal := findMinimum(firstMap[firstVal]+1, countOnSecond)
					if newVal > previousVal {
						totalOverlap++
					}
				}

				firstMap[firstVal]++
				firstTotalKeys++
			} else {
				firstDone = true
			}
		case secondVal, more := <-second:
			if more {
				if _, ok := secondMap[secondVal]; !ok {
					secondMap[secondVal] = 0
					secondDistinctKeys++

					if _, ok := firstMap[secondVal]; ok {
						distinctOverlap++
					}
				}

				if countOnFirst, ok := firstMap[secondVal]; ok {
					previousVal := findMinimum(secondMap[secondVal], countOnFirst)
					newVal := findMinimum(secondMap[secondVal]+1, countOnFirst)
					if newVal > previousVal {
						totalOverlap++
					}
				}

				secondMap[secondVal]++
				secondTotalKeys++
			} else {
				secondDone = true
			}
		}
	}

	result := &IntersectionResult{
		First: &FileResult{
			KeyCount:         firstTotalKeys,
			DistinctKeyCount: firstDistinctKeys,
		},
		Second: &FileResult{
			KeyCount:         secondTotalKeys,
			DistinctKeyCount: secondDistinctKeys,
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

func findOverlapsParallel(firstKeys, secondKeys map[string]int) (int, int) {
	count := 0
	totalOverlap := 0

	wg := sync.WaitGroup{}
	wg.Add(len(firstKeys))

	for fk, fv := range firstKeys {
		go func(fk string, fv int) {
			defer wg.Done()

			if sv, ok := secondKeys[fk]; ok {
				count++

				if fv < sv {
					totalOverlap += fv
				} else {
					totalOverlap += sv
				}
			}
		}(fk, fv)
	}

	wg.Wait()

	return count, totalOverlap
}

// IntersectionResult represents result of intersection count
type IntersectionResult struct {
	First           *FileResult
	Second          *FileResult
	TotalOverlap    int
	DistinctOverlap int
}

// FileResult represents result of a file key count
type FileResult struct {
	KeyCount         int
	DistinctKeyCount int
}

func findMinimum(a, b int) int {
	if a < b {
		return a
	}

	return b
}
