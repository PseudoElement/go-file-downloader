run-dev:
	go run main.go

run-build:
	go build && ./go-file-downloader

test:
	go test ./... -v