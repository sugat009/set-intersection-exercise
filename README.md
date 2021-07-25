# set-intersection

## To build

Move this folder to your `$GOPATH`

If GOMODULE is not enabled then set
```
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

## To run
```sh
go run main.go --first-file=[path_to_first-file]] --second-file=[oath_to_second_file] --key=foo
```