build:
	CGO_ENABLED=1 go build -o ./build/betterreads

test:
	DB_TYPE=sqlite DATABASE_URL=:memory: go test ./...

fmt:
	go fmt ./...

fmt-check:
	go fmt -n ./...

run:
	DB_TYPE=sqlite DATABASE_URL=local.db go run main.go

vet:
	go vet ./...

.PHONY: build