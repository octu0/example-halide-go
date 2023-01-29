.PHONY: all
all: generate test

.PHONY: generate
generate:
	go generate .

.PHONY: test
test:
	go test -v .

