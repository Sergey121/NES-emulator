test:
	go test -v ./...

test-ppu:
	go test -v ./internal/ppu/...

test-cpu:
	go test -v ./internal/cpu/...

test-rom:
	go test -v ./internal/rom/...

run:
	go run cmd/main.go

build:
	go build -o bin/main cmd/main.go

