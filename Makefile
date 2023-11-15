build:
	go build -o encrypted-messages cmd/main.go -config=config.yaml.example

run:
	go run cmd/main.go -config=config.yaml.example