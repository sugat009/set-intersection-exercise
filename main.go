package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/pterm/pterm"
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
		Usage: "Given two input files in CSV format and a key, the program outputs the total no. of keys and distinct no. of keys in each file. It also provides the total overlap and distinct overlap between the two files.",
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
	startedAt := time.Now()

	firstKeys := make(chan string, cfg.BufferSize)
	secondKeys := make(chan string, cfg.BufferSize)

	errorCh := make(chan error)

	taskSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Reading files %s and %s using key: %s", cfg.FirstSource, cfg.SecondSource, cfg.Key))

	var readFirst, readSecond bool
	go func() {
		if err := fileToKeysChannel(cfg.FirstSource, cfg.Key, firstKeys); err != nil {
			errorCh <- err
		}
		pterm.Success.Println("Processed first file")
		taskSpinner.UpdateText("Completed reading first file. Still reading second file...")
		readFirst = true
		if readSecond {
			taskSpinner.UpdateText("Finding intersections...")
		}
	}()

	go func() {
		if err := fileToKeysChannel(cfg.SecondSource, cfg.Key, secondKeys); err != nil {
			errorCh <- err
		}
		pterm.Success.Println("Processed second file")
		taskSpinner.UpdateText("Completed reading second file. Still reading first file...")
		readSecond = true
		if readFirst {
			taskSpinner.UpdateText("Finding intersections...")
		}
	}()

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
		taskSpinner.Fail("Process crashed")
		return err
	case result := <-resultCh:
		taskSpinner.Success(fmt.Sprintf("Process completed. Elapsed: %s", time.Since(startedAt).String()))
		showResult(cfg.FirstSource, cfg.SecondSource, result)
	}
	return nil
}

func fileToKeysChannel(filePath, key string, output chan<- string) error {
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

	if err := reader.ReadKeysFromCsvIntoChannel(key, bufio.NewReader(file), output); err != nil {
		return errors.Wrapf(err, "while processing file: %s", filePath)
	}

	return nil
}

func showResult(firstFilePath, secondFilePath string, result counter.IntersectionResult) {
	err := pterm.DefaultTable.WithHasHeader().WithData(pterm.TableData{
		{
			"Total keys in first table",
			"Distinct keys in first table",
			"Total keys in second table",
			"Distinct keys in second table",
			"Total Overlap",
			"Distinct Overlap",
		},
		{
			fmt.Sprintf("%v", result.First.KeyCount),
			fmt.Sprintf("%v", result.First.DistinctKeyCount),
			fmt.Sprintf("%v", result.Second.KeyCount),
			fmt.Sprintf("%v", result.Second.DistinctKeyCount),
			fmt.Sprintf("%v", result.TotalOverlap),
			fmt.Sprintf("%v", result.DistinctOverlap),
		},
	}).Render()
	if err != nil {
		log.Error(err.Error())
	}
}
