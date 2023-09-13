init:
	git config core.hooksPath git-hooks

pre-commit:
	cd ./runtime && go generate
	go fmt && go mod tidy
	$(MAKE) integrations-generate
	$(MAKE) build-readme

integrations-generate:
	cd ./cmd/internal/integration && go generate

integrations-run:
	go run ./cmd/internal/integration/generated_xgen_cli.go

build-readme:
	go run ./cmd/docbuilder/main.go build -t README.gomd -o README.md