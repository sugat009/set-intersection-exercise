package app

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/tav/golly/log"

	"github.com/rickyshrestha/infosum-interview-exercise/internal/counter"
)

// ReadKeyFromFileFunc signature of function that can be used to read keys from a file
type ReadKeyFromFileFunc func(key string, reader io.Reader, keysOuput chan<- string) error

// NewApp creates a new app for finding set intersection using the func passed in the parameter to parse keys from the input files
func NewApp(readKeysFunc ReadKeyFromFileFunc) App {
	return App{
		readKeyFromFile: readKeysFunc,
	}
}

// App represents a run of app
type App struct {
	// func to use to parse input files for keys
	readKeyFromFile ReadKeyFromFileFunc
}

// RuntimeParam are parameters used for running the app
type RuntimeParam struct {
	FirstSource, SecondSource, Key string
	BufferSize                     int
}

// Start starts the read from file and processing the intersections
func (a *App) Start(param RuntimeParam) (counter.IntersectionResult, error) {
	if a.readKeyFromFile == nil {
		return counter.IntersectionResult{}, errors.New("function to parse input files for keys is not set")
	}

	// read first tile
	firstKeys := make(chan string, param.BufferSize)
	secondKeys := make(chan string, param.BufferSize)

	errorCh := make(chan error)

	go func() {
		if err := a.readFileIntoKeysChannel(param.FirstSource, param.Key, firstKeys); err != nil {
			errorCh <- err
		}
	}()

	// read second file
	go func() {
		if err := a.readFileIntoKeysChannel(param.SecondSource, param.Key, secondKeys); err != nil {
			errorCh <- err
		}
	}()

	// find overlaps
	resultCh := make(chan counter.IntersectionResult)
	go func() {
		result, err := counter.FindSetIntersection(firstKeys, secondKeys)
		if err != nil {
			errorCh <- errors.Wrap(err, "while finding intersection")
		}
		resultCh <- result
	}()

	select {
	case err := <-errorCh:
		return counter.IntersectionResult{}, err
	case result := <-resultCh:
		return result, nil
	}
}

func (a *App) readFileIntoKeysChannel(filePath, key string, output chan<- string) error {
	defer close(output)

	file, err := os.Open(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to read file: %s", filePath)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("unable to close file: %s", filePath)
		}
	}()

	if err := a.readKeyFromFile(key, file, output); err != nil {
		return errors.Wrapf(err, "while processing file: %s", filePath)
	}

	return nil
}
