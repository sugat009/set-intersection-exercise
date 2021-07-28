package app

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/rickyshrestha/set-intersection-exercise/internal/counter"
)

func mockReadKeyFromFile(key string, reader io.Reader, keysOutput chan<- string) error {
	if key != "key" {
		return errors.New("invalid test key for mock")
	}

	r := bufio.NewReader(reader)
	line, _, _ := r.ReadLine()

	for _, l := range strings.Split(string(line), ",") {
		keysOutput <- l
	}

	return nil
}

func Test_Start_Success(t *testing.T) {
	a := NewApp(mockReadKeyFromFile)
	res, err := a.Start(RuntimeParam{
		FirstSource:  "./testdata/first.txt",
		SecondSource: "./testdata/second.txt",
		Key:          "key",
		BufferSize:   64,
	})

	assert.NoError(t, err)
	assert.Equal(t, counter.IntersectionResult{
		First: counter.FileResult{
			KeyCount:         8,
			DistinctKeyCount: 6,
		},
		Second: counter.FileResult{
			KeyCount:         9,
			DistinctKeyCount: 6,
		},
		TotalOverlap:    11,
		DistinctOverlap: 4,
	}, res)
}
