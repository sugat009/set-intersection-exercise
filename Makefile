build:
	go mod tidy
	go build -o set-intersection-exercise main.go

test:
	go test -v -cover -race ./...

lint:
	golangci-lint run
	golint ./...

RUNFLAGS = 
RUNFLAGS += "--first-file=./testdata/tiny_1.csv"
RUNFLAGS += "--second-file=./testdata/tiny_2.csv"
RUNFLAGS += "--key=foo"

run: 
	go mod tidy
	go run main.go $(RUNFLAGS)

run-med: 
	go run main.go --first-file=./testdata/med_1.csv --second-file=./testdata/med_2.csv --key=foo


run-large: 
	go run main.go --first-file=./testdata/large_1.csv --second-file=./testdata/large_2.csv --key=foo

run-x-large: 
	go run main.go --first-file=./testdata/x_large_1.csv --second-file=./testdata/x_large_2.csv --key=foo

bench:
	cd internal/counter && go test -benchmem -bench=.
	cd internal/reader && go test -benchmem -bench=.

check: test lint
	astitodo . && go mod tidy && gofumpt -s -w . && goimports -w . && gocyclo -top 5 . && errcheck ./... && goconst ./... && go vet ./... && go clean -testcache ./... 