.PHONY: deps
deps:
	go mod tidy

.PHONY: build
build:
	go build \
	-o ./artifacts/svc \
	.