tidy:
	go mod tidy
build:
	go build -o encrypted-messages cmd/main.go

run:
	go run cmd/main.go -config=config.yaml