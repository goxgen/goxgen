# Initialize the git hooks
init:
	git config core.hooksPath git-hooks

# Run the pre-commit hook
pre-commit:
	$(MAKE) build

# Generate the integration tests
integrations-generate:
	cd ./cmd/internal/integration && go generate

# Generate the runtime files that are used by goxgen to generate the code
runtime-generate:
	cd ./runtime && go generate

# Run the integration tests
integrations-run:
	cd ./cmd/internal/integration/; go run generated_xgen_cli.go run --graphql-playground-enabled=true

# Build the README.md file from the README.gomd file
build-readme:
	go run ./cmd/docbuilder/main.go build -t README.gomd -o README.md

# Build all the things
build:
	cd ./runtime && go generate
	go fmt && go mod tidy
	$(MAKE) integrations-generate
	$(MAKE) build-readme

# Run the tests including the integration tests
test:
	go test -v ./...
	cd ./cmd/internal/integration/; go run generated_xgen_cli.go run --test=true

