BINARY := ,ht
BIN_DIR := bin
GOCACHE ?= $(CURDIR)/.cache/go-build

.PHONY: build test clean run

build:
	mkdir -p $(BIN_DIR)
	GOCACHE=$(GOCACHE) go build -o $(BIN_DIR)/$(BINARY) .

test:
	GOCACHE=$(GOCACHE) go test ./...

clean:
	rm -rf $(BIN_DIR) .cache

run:
	GOCACHE=$(GOCACHE) go run . $(ARGS)
