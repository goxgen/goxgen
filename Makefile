pre-commit:
	cd ./runtime && go generate
	go fmt && go mod tidy
	$(MAKE) integrations-generate

integrations-generate:
	cd ./internal/integrations && go generate