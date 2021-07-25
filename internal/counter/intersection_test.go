package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	bufferSize = 32
	keyLength  = 500
	rowCount   = 500000
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
	assert.Equal(t, 4, total)
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
	assert.Equal(t, &IntersectionResult{
		First: &FileResult{
			KeyCount:         8,
			DistinctKeyCount: 6,
		},
		Second: &FileResult{
			KeyCount:         9,
			DistinctKeyCount: 6,
		},
		DistinctOverlap: 4,
		TotalOverlap:    5,
	}, res)
}

func Benchmark_findOverlaps(b *testing.B) {
	for i := 0; i < b.N; i++ {

		input1 := make(map[string]int)
		for i := 0; i < 100000; i++ {
			input1[getRandomString(5)] = i
		}

		input2 := make(map[string]int)
		for i := 0; i < 100000; i++ {
			input2[getRandomString(5)] = i
		}

		_, _ = findOverlaps(input1, input2)
	}
}

func Benchmark_findOverlaps_Parallel(b *testing.B) {
	for i := 0; i < b.N; i++ {

		input1 := make(map[string]int)
		for i := 0; i < 100000; i++ {
			input1[getRandomString(5)] = i
		}

		input2 := make(map[string]int)
		for i := 0; i < 100000; i++ {
			input2[getRandomString(5)] = i
		}

		_, _ = findOverlaps_Parallel(input1, input2)
	}
}

func Benchmark_findSetIntersection_v1(b *testing.B) {

	for i := 0; i < b.N; i++ {

		first := make(chan string, bufferSize)
		second := make(chan string, bufferSize)

		//add to the channel
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

		_, _ = findSetIntersection_v1(first, second)
	}
}

func Benchmark_findSetIntersection_v2(b *testing.B) {

	for i := 0; i < b.N; i++ {

		first := make(chan string, bufferSize)
		second := make(chan string, bufferSize)

		//add to the channel
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

		_, _ = findSetIntersection_v2(first, second)
	}
}

func Benchmark_findSetIntersection_v3(b *testing.B) {

	for i := 0; i < b.N; i++ {

		first := make(chan string, bufferSize)
		second := make(chan string, bufferSize)

		//add to the channel
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

		_, _ = findSetIntersection_v3(first, second)
	}
}

func Benchmark_findSetIntersection_v4(b *testing.B) {

	for i := 0; i < b.N; i++ {

		first := make(chan string, bufferSize)
		second := make(chan string, bufferSize)

		//add to the channel
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

		_, _ = findSetIntersection_v4(first, second)
	}
}

func populateChannel(input chan<- string) {
	for i := 0; i < rowCount; i++ {
		input <- getRandomString(keyLength)
	}
}
