# set-intersection-exercise

Given two input files in CSV format and a key, the program outputs the total no. of keys and distinct no. of keys in each file. It also provides the total overlap and distinct overlap between the two files.

## Installing

You can download the pre-built binaries from the [release page](https://github.com/rickyshrestha/set-intersection-exercise/releases)

### Building from source

You can clone this repository into your go workspace and build from source.

It uses `GOMODULES`. If it is not already `on` please do so using:

```sh
export GO111MODULE="on"
```

Now fetch dependent modules and build the executable

```sh
go mod tidy && go build -o set-intersection-exercise main.go
```

## To Run

```sh
./set-intersection-exercise --first-file=[path_to_first_file]] --second-file=[oath_to_second_file] --key=foo
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

### Need Help ?

```sh
./set-intersection-exercise --help
```

## To Test

```sh
go test -v -cover -race ./...
```

## Makefile

If you have `make` installed, then you can use it to lint, run, test and benchmark the source code. Please see `Makefile` for more info.
