# set-intersection

Given two input files in CSV format and a key, the program outputs the total no. of keys and distinct no. of keys in each file. It also provides the total overlap and distinct overlap between the two files.

## To run

```sh
go run main.go --first-file=[path_to_first_file]] --second-file=[oath_to_second_file] --key=foo
```

### Output

```text
Count of keys ([path_to_first_file]]):           [no. of keys]
Count of distinct keys ([path_to_first_file]):   [no. of distinct keys]
Count of keys ([path_to_second_file]):           [no. of keys]
Count of distinct keys ([path_to_second_file]):  [no. of distinct keys]
Total overlapping keys:                          [total overlapping keys]
Distinct overlapping keys:                       [total distinct keys]
```

## Help

```sh
go run main.go --help
```

```text
NAME:
   set-intersection - Given two input files in CSV format and a key, the program outputs the total no. of keys and distinct no. of keys in each file. It also provides the total overlap and distinct overlap between the two files.

USAGE:
    [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --first-file value   path to the first of the two files to compare [$FIRST_FILE]
   --second-file value  path to the second of the two files to compare [$SECOND_FILE]
   --key value          column in the csv file to be used as the key for comparison [$KEY]
   --buffer-size value  buffer size for no. of records to load from file to process (default: 64) [$BUFFER_SIZE]
   --help, -h           show help
```

## To build

Move this folder to your `$GOPATH`

If GOMODULE is not enabled then set

```sh
export GO111MODULE="on"
```

```sh
go get ./...
go mod tidy
go build -o set-intersection 
```

## To Test

```sh
go test -v -cover -race ./...
```

## Makefile

If you have `make` installed, then you can use it to lint, run, test and benchmark the source code. Please see `Makefile` for more info.
