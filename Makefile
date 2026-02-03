MODULE := github.com/17882237881/MaaS

.PHONY: help	help tidy test vet

help:
	@echo "Targets: tidy, test, vet"

tidy:
	go mod tidy

test:
	go test ./...

vet:
	go vet ./...
