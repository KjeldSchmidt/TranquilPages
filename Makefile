build:
	go build -o ./build/tranquil-pages

test:
	go test ./...

fmt:
	go fmt ./...
	cd infra && terraform fmt -recursive .

fmt-check:
	@echo "Checking Go code formatting..."
	@if gofmt -l . | tee /dev/stderr | grep .; then \
		echo "Go files are not formatted!"; \
		exit 1; \
	fi
	@echo "Checking Terraform formatting..."
	cd infra && terraform fmt -recursive -check .

run:
	sudo docker start mongodb
	DB_URL="mongodb://localhost:27017" go run main.go

vet:
	go vet ./...

tf-validate:
	cd infra/base && terraform init && terraform validate
	cd infra/env/dev && terraform init && terraform validate

tf-apply-auto:
	cd "infra/env/${env}" && terraform init && terraform apply -auto-approve

quality-gates: fmt vet test build tf-validate
	echo "✅✅✅"

build-image:
	sudo docker build . -t kjeldschmidt2/tranquil-pages:latest

push-image:
	sudo docker push kjeldschmidt2/tranquil-pages:latest

reset-local-db:
	sudo docker stop mongodb
	sudo docker rm mongodb
	sudo docker run -d --name mongodb -p 27017:27017 mongo:latest

.PHONY: build