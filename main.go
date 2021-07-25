package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rickyshrestha/infosum-interview-exercise/internal/counter"
	"github.com/rickyshrestha/infosum-interview-exercise/internal/reader"
	"github.com/tav/golly/log"
	"github.com/urfave/cli"
)

const (
	flagFirstFile  = "first-file"
	flagSecondFile = "second-file"
	flagKey        = "key"
	flagBufferSize = "buffer-size"
)

func main() {

	app := &cli.App{
		Name:  "set-intersection",
		Usage: "find intersections beetween two csv files dataset using a common key",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   flagFirstFile,
				EnvVar: "FIRST_FILE",
				Usage:  "path to the first of the two files to compare",
			},
			cli.StringFlag{
				Name:   flagSecondFile,
				EnvVar: "SECOND_FILE",
				Usage:  "path to the second of the two files to compare",
			},
			cli.StringFlag{
				Name:   flagKey,
				EnvVar: "KEY",
				Usage:  "column in the csv file to be used as the key for comparison",
			},
			cli.IntFlag{
				Name:   flagBufferSize,
				EnvVar: "BUFFER_SIZE",
				Usage:  "buffer size for no. of records to load from file to process",
				Value:  64,
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(context *cli.Context) error {
	cfg, err := parseAppConfig(context)
	if err != nil {
		return errors.Wrap(err, "invalid application configs")
	}

	if err := startApplication(cfg); err != nil {
		return errors.Wrap(err, "while running application")
	}

	return nil
}

func parseAppConfig(context *cli.Context) (appConfig, error) {

	config := appConfig{}
	config.BufferSize = context.Int(flagBufferSize)

	if config.BufferSize <= 0 {
		return config, errors.Errorf("invalid buffer size (%s): %v", flagBufferSize, config.BufferSize)
	}

	config.FirstSource = context.String(flagFirstFile)
	if config.FirstSource == "" {
		return config, errors.New("first source file is empty")
	}

	config.SecondSource = context.String(flagSecondFile)
	if config.SecondSource == "" {
		return config, errors.New("second source file is empty")
	}

	config.Key = context.String(flagKey)

	return config, nil
}

type appConfig struct {
	FirstSource, SecondSource, Key string
	BufferSize                     int
}

func startApplication(cfg appConfig) error {

	benchStart := time.Now()

	firstKeys := make(chan string, cfg.BufferSize)
	secondKeys := make(chan string, cfg.BufferSize)

	errorCh := make(chan error)

	go func() {
		if err := fileToKeysChannel(cfg.FirstSource, cfg.Key, firstKeys); err != nil {
			errorCh <- err
		}
		log.Infof("Processed: %s", cfg.FirstSource)
	}()

	go func() {
		if err := fileToKeysChannel(cfg.SecondSource, cfg.Key, secondKeys); err != nil {
			errorCh <- err
		}
		log.Infof("Processed: %s", cfg.SecondSource)
	}()

	resultCh := make(chan *counter.IntersectionResult)
	go func() {
		log.Infof("Finding intersections in: %s & %s using key: %s ...", cfg.FirstSource, cfg.SecondSource, cfg.Key)
		result, err := counter.FindSetIntersection(firstKeys, secondKeys)
		if err != nil {
			errorCh <- errors.Wrap(err, "while finding intersection")
		}
		resultCh <- result
	}()

	select {
	case err := <-errorCh:
		return err
	case result := <-resultCh:
		showResult(cfg.FirstSource, cfg.SecondSource, result)
		log.Infof("elapsed: %s", time.Since(benchStart).String())
	}
	return nil
}

func fileToKeysChannel(filePath, key string, output chan<- string) error {

	defer close(output)

	file, err := os.Open(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to read file: %s", filePath)
	}
	defer file.Close()

	if err := reader.ReadKeysFromCsvIntoChannel(key, bufio.NewReader(file), output); err != nil {
		return errors.Wrapf(err, "while processing file: %s", filePath)
	}

	return nil
}

func showResult(firstFilePath, secondFilePath string, result *counter.IntersectionResult) {

	fmt.Printf("Count of keys (%s):\t\t%v\n", firstFilePath, result.First.KeyCount)
	fmt.Printf("Count of distinct keys (%s):\t%v\n", firstFilePath, result.First.DistinctKeyCount)
	fmt.Printf("Count of keys (%s):\t\t%v\n", secondFilePath, result.Second.KeyCount)
	fmt.Printf("Count of distinct keys (%s):\t%v\n", secondFilePath, result.Second.DistinctKeyCount)

	fmt.Printf("Total overlapping keys:\t\t\t\t%v\n", result.TotalOverlap)
	fmt.Printf("Distinct overlapping keys:\t\t\t%v\n", result.DistinctOverlap)
}
