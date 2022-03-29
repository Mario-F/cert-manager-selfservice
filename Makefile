SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=-v

# Run generators for automatic code generation
generate:
	go generate ./...

# Run all the tests
test: generate
	LC_ALL=C go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=5m

# Run all the linters, not working with 1.18 because of https://github.com/golangci/golangci-lint/issues/2649
lint: generate
	golangci-lint run ./...
	misspell -error **/*

# Run gorelease dry
godry:
	goreleaser --snapshot --skip-publish --rm-dist

# Run a build process
build:
	go build -o ./build/server ./cmd/app/*.go
