pre-commit:
	cd ./runtime && go generate
	go fmt && go mod tidy