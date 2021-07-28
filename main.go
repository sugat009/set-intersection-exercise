package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/tav/golly/log"
	"github.com/urfave/cli"

	"github.com/rickyshrestha/set-intersection-exercise/internal/app"
	"github.com/rickyshrestha/set-intersection-exercise/internal/counter"
	"github.com/rickyshrestha/set-intersection-exercise/internal/reader"
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
	startedAt := time.Now()

	cfg, err := parseAppConfig(context)
	if err != nil {
		return errors.Wrap(err, "invalid application configs")
	}

	counterApp := app.NewApp(reader.ReadKeysFromCsvIntoChannel)

	result, err := counterApp.Start(cfg)
	if err != nil {
		return errors.Wrap(err, "while running application")
	}

	showResult(cfg.FirstSource, cfg.SecondSource, result)
	pterm.DefaultSpinner.Success(fmt.Sprintf("Process completed. Elapsed: %s", time.Since(startedAt).String()))
	return nil
}

func parseAppConfig(context *cli.Context) (app.RuntimeParam, error) {
	config := app.RuntimeParam{}
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
