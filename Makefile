GO = $(shell which go)

semsearch: $(shell find . -name '*.go')
	$(GO) build ./cmd/semsearch
