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

## Help

```text
./set-intersection --help
```
