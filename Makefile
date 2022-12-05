.PHONY: run/todo
run/todo:
	@go run ./cmd/todo

.PHONY: build
build:
	@go build ./cmd/todo