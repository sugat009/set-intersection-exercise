SERVICE=set-intersection-exercise

build:
	go mod tidy
	go build -o $(SERVICE) main.go

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
	cd internal/app && go test -benchmem -bench=.

check: test lint
	astitodo . && go mod tidy && gofumpt -s -w . && goimports -w . && gocyclo -top 5 . && errcheck ./... && goconst ./... && go vet ./... && go clean -testcache ./... 

release:
	GOOS=darwin GOARCH=amd64  go build -o releases/$(SERVICE)-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 go build -o releases/$(SERVICE)-linux-amd64 main.go
	GOOS=linux GOARCH=386 go build -o releases/$(SERVICE)-linux-386 main.go
	GOOS=windows GOARCH=amd64 go build -o releases/$(SERVICE)-windows-amd64.exe main.go
	GOOS=windows GOARCH=386 go build -o releases/$(SERVICE)-windows-386.exe main.go