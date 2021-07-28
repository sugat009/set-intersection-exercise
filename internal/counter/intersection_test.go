package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	bufferSize = 32
	keyLength  = 500
	rowCount   = 50000
)

func Test_findOverlaps_Empty(t *testing.T) {
	distinct, total := findOverlaps(map[string]int{}, map[string]int{})
	assert.Equal(t, 0, distinct)
	assert.Equal(t, 0, total)
}

func Test_findOverlaps(t *testing.T) {
	distinct, total := findOverlaps(map[string]int{
		"a": 1,
	}, map[string]int{
		"a": 1,
	})
	assert.Equal(t, 1, distinct)
	assert.Equal(t, 1, total)
}

func Test_findOverlaps_NoOverlap(t *testing.T) {
	distinct, total := findOverlaps(map[string]int{
		"a": 1,
	}, map[string]int{
		"b": 1,
	})
	assert.Equal(t, 0, distinct)
	assert.Equal(t, 0, total)
}

func Test_findOverlaps_MultipleOverlaps(t *testing.T) {
	distinct, total := findOverlaps(map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}, map[string]int{
		"a": 3,
		"b": 2,
		"c": 1,
		"e": 2,
	})
	assert.Equal(t, 3, distinct)
	assert.Equal(t, 10, total)
}

func Test_FindSetIntersection(t *testing.T) {
	first := make(chan string, bufferSize)
	second := make(chan string, bufferSize)

	go func() {
		defer close(first)
		first <- "a"
		first <- "b"
		first <- "c"
		first <- "d"
		first <- "d"
		first <- "e"
		first <- "f"
		first <- "f"
	}()

	go func() {
		defer close(second)
		second <- "a"
		second <- "c"
		second <- "c"
		second <- "d"
		second <- "f"
		second <- "f"
		second <- "f"
		second <- "x"
		second <- "y"
	}()

	res, err := FindSetIntersection(first, second)
	assert.NoError(t, err)
	assert.Equal(t, IntersectionResult{
		First: FileResult{
			KeyCount:         8,
			DistinctKeyCount: 6,
		},
		Second: FileResult{
			KeyCount:         9,
			DistinctKeyCount: 6,
		},
		DistinctOverlap: 4,
		TotalOverlap:    11,
	}, res)
}

func Test_FindSetIntersection_Nil(t *testing.T) {
	_, err := FindSetIntersection(nil, nil)
	assert.Error(t, err)
}

func Test_FindSetIntersection_Empty(t *testing.T) {
	first := make(chan string, bufferSize)
	second := make(chan string, bufferSize)

	go close(first)
	go close(second)

	res, err := FindSetIntersection(first, second)
	assert.NoError(t, err)
	assert.Equal(t, IntersectionResult{
		First: FileResult{
			KeyCount:         0,
			DistinctKeyCount: 0,
		},
		Second: FileResult{
			KeyCount:         0,
			DistinctKeyCount: 0,
		},
		DistinctOverlap: 0,
		TotalOverlap:    0,
	}, res)
}

func Test_FindSetIntersection_ClosedChannel(t *testing.T) {
	first := make(chan string, bufferSize)
	second := make(chan string, bufferSize)
	close(first)
	close(second)

	res, err := FindSetIntersection(first, second)
	assert.NoError(t, err)
	assert.Equal(t, IntersectionResult{
		First: FileResult{
			KeyCount:         0,
			DistinctKeyCount: 0,
		},
		Second: FileResult{
			KeyCount:         0,
			DistinctKeyCount: 0,
		},
		DistinctOverlap: 0,
		TotalOverlap:    0,
	}, res)
}

func Benchmark_findOverlaps(b *testing.B) {
	for i := 0; i < b.N; i++ {

		input1 := make(map[string]int)
		for i := 0; i < 5000; i++ {
			input1[getRandomString(3)] = i
		}

		input2 := make(map[string]int)
		for i := 0; i < 5000; i++ {
			input2[getRandomString(3)] = i
		}

		_, _ = findOverlaps(input1, input2)
	}
}

func Benchmark_FindSetIntersection(b *testing.B) {
	for i := 0; i < b.N; i++ {

		first := make(chan string, bufferSize)
		second := make(chan string, bufferSize)

		// add to the channel
		go func() {
			defer close(first)
			for i := 0; i < rowCount; i++ {
				first <- getRandomString(keyLength)
			}
		}()

		go func() {
			defer close(second)
			for i := 0; i < rowCount; i++ {
				second <- getRandomString(keyLength)
			}
		}()

		_, _ = FindSetIntersection(first, second)
	}
}
