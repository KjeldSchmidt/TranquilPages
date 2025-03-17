build:
	cd backend && make build

test:
	cd backend && make test

fmt:
	cd backend && make fmt
	cd infra && terraform fmt -recursive .

fmt-check:
	cd backend && make fmt-check
	@echo "Checking Terraform formatting..."
	cd infra && terraform fmt -recursive -check .

run:
	sudo docker start mongodb
	cd backend && DB_URL="mongodb://localhost:27017" go run main.go

lint:
	cd backend && make vet
	make tf-validate

tf-validate:
	cd infra/base && terraform init && terraform validate
	cd infra/env/dev && terraform init && terraform validate

tf-apply-auto:
	@if [ -z "$(env)" ]; then echo "Error: env is not set. Please pass by name: env=<dev|staging|prod>."; exit 1; fi
	cd "infra/env/${env}" && terraform init && terraform apply -auto-approve

quality-gates: fmt lint test build tf-validate
	echo "✅✅✅"

build-image:
	cd backend && make build-image

push-image:
	cd backend && make push-image

reset-local-db:
	sudo docker stop mongodb
	sudo docker rm mongodb
	sudo docker run -d --name mongodb -p 27017:27017 mongo:latest

install-hooks:
	git config core.hooksPath ./.githooks

.PHONY: build