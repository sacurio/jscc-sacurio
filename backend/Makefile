BINARY_NAME=chat-bot

# clean previous unused dependencies and install the provided ones.
install: clean-dependencies
	go get ./...
	go get -u gorm.io/gorm

# clean previous unused dependencies.
clean-dependencies:
	go mod tidy

# run linter
lint:
	golangci-lint run ./...

# start the app.
start:
	go run cmd/main.go

# clean bin folder.
clean:
	go clean
	rm bin/*

# clean bind folder and build the app.
build: clean
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin cmd/main.go
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux cmd/main.go
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows cmd/main.go

# build and run chatbot app.
start: build
	./bin/${BINARY_NAME}-linux

# build and run chatbot app.
run:
	go run cmd/main.go

# executes docker-compose file.
docker-up:
	docker-compose up -d

# stops the docker-compose services.
docker-stop:
	docker-compose down