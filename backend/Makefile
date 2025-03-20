build:
	go mod download
	go build -o ./build/tranquil-pages

install-deps:
	go install gotest.tools/gotestsum@latest

test:
	gotestsum --format testname ./...

fmt:
	go fmt ./...

fmt-check:
	@echo "Checking Go code formatting..."
	@if gofmt -l . | tee /dev/stderr | grep .; then \
		echo "Go files are not formatted!"; \
		exit 1; \
	fi

run:
	 dotenvx run -f ../.env.local -- go run main.go

lint:
	go vet ./...
	go mod tidy

quality-gates: fmt lint test build
	echo "✅✅✅"

build-image:
	docker build . -t kjeldschmidt2/tranquil-pages:latest

push-image:
	docker push kjeldschmidt2/tranquil-pages:latest

reset-local-db:
	sudo docker stop mongodb
	sudo docker rm mongodb
	sudo docker run -d --name mongodb -p 27017:27017 mongo:latest

start-local-db:
	sudo docker start mongodb

autostart-local-db:
	sudo docker update --restart unless-stopped mongodb

.PHONY: build install-deps