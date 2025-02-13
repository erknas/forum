build:
	@go build -o bin/forum cmd/main.go

run: build
	@./bin/forum

generate:
	go run github.com/99designs/gqlgen generate 