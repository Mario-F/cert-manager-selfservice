SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=-v

# Run generators for automatic code generation
generate:
	go generate ./...
	oapi-codegen -generate types -package api openapi.yaml > internal/gen/api/types.go
	oapi-codegen -generate client -package api openapi.yaml > internal/gen/api/client.go
	oapi-codegen -generate server -package api openapi.yaml > internal/gen/api/server.go
	oapi-codegen -generate spec -package api openapi.yaml > internal/gen/api/spec.go
	openapi-generator generate -i openapi.yaml -g typescript-axios -o web/src/gen/api

# Run all the tests
test:
	LC_ALL=C go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=5m

# Run all the linters, not working with 1.18 because of https://github.com/golangci/golangci-lint/issues/2649
lint:
	golangci-lint run ./...
	misspell -error **/*

# Run gorelease dry
godry:
	goreleaser --snapshot --skip-publish --rm-dist

# Run a build process
build:
	go build -o ./build/server ./cmd/app/*.go
