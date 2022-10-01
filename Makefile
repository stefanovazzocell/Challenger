# Testing
.PHONY: test
test:
	@go test -cover -race ./provider
	@go test -cover -race ./provider/utils

.PHONY: bench
bench:
	@go test -cover -benchmem -cover -bench . ./provider
	@go test -cover -benchmem -cover -bench . ./provider/utils

# Debugging
.PHONY: demo
demo:
	@go run -race ./demo