BINARY_NAME=aws-profile
INSTALL_PATH=$(shell go env GOPATH)/bin
DIST=dist

.PHONY: install build run deps clean release

install: build
	@mkdir -p $(INSTALL_PATH)
	@cp $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "Instalado en $(INSTALL_PATH)/$(BINARY_NAME)"

build: deps
	@go build -o $(BINARY_NAME) .

run: deps
	@go run .

deps:
	@go mod download
	@go mod tidy

clean:
	@rm -f $(BINARY_NAME)
	@rm -rf $(DIST)

release: deps
	@mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o $(DIST)/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o $(DIST)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(DIST)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o $(DIST)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Binarios en $(DIST)/"
	@ls -la $(DIST)/