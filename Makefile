build:
	CGO_ENABLED=1 go build -o ./build/betterreads

test:
	DB_TYPE=sqlite DATABASE_URL=:memory: go test ./...

fmt:
	go fmt ./...

fmt-check:
	@echo "Checking Go code formatting..."
	@if gofmt -l . | tee /dev/stderr | grep .; then \
		echo "Go files are not formatted!"; \
		exit 1; \
	fi
	@echo "Formatting passes!"

run:
	DB_TYPE=sqlite DATABASE_URL=local.db go run main.go

vet:
	go vet ./...

.PHONY: build