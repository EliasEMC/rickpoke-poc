run:             ## arranca el servidor
	go run ./cmd/server

build:           ## compila binario
	go build -o bin/server ./cmd/server

test:            ## tests con cobertura
	go test ./... -cover

lint:            ## ejecuta golangci-lint
	golangci-lint run

fmt:             ## formatea
	go fmt ./...

tidy:            ## resuelve dependencias
	go mod tidy
