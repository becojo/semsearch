GO := $(shell which go)
TARGET := semsearch
SRC := $(shell find pkg cmd -name '*.go' -not -path '*_test.go')

semsearch: $(SRC)
	$(GO) build -o $(TARGET) ./cmd/semsearch

clean:
	rm -f $(TARGET) coverage.out coverage.html

test:
	$(GO) test -v ./pkg/...

coverage:
	go test ./pkg/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html


format: format.go format.keep-sorted

format.go:
	$(GO) fmt ./pkg/...
	$(GO) fmt ./cmd/...

format.keep-sorted:
	keep-sorted $(SRC)
