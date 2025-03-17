build:
	go build -o ./build/tranquil-pages

test:
	go test ./...

fmt:
	go fmt ./...

fmt-check:
	@echo "Checking Go code formatting..."
	@if gofmt -l . | tee /dev/stderr | grep .; then \
		echo "Go files are not formatted!"; \
		exit 1; \
	fi

run:
	sudo docker start mongodb
	DB_URL="mongodb://localhost:27017" go run main.go

vet:
	go vet ./...

quality-gates: fmt vet test build
	echo "✅✅✅"

build-image:
	docker build . -t kjeldschmidt2/tranquil-pages:latest

push-image:
	docker push kjeldschmidt2/tranquil-pages:latest

reset-local-db:
	sudo docker stop mongodb
	sudo docker rm mongodb
	sudo docker run -d --name mongodb -p 27017:27017 mongo:latest

.PHONY: build