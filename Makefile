BINARY_NAME=aws-profile
INSTALL_PATH=$(shell go env GOPATH)/bin

.PHONY: install build run deps clean

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