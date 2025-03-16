build:
	CGO_ENABLED=1 go build -o ./build/tranquil-pages

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
	DB_TYPE=sqlite DATABASE_URL=local.db go run main.go

vet:
	go vet ./...

tf-validate:
	cd infra/base && terraform init && terraform validate
	cd infra/env/dev && terraform init && terraform validate

tf-apply-auto:
	cd "infra/env/${env}" && terraform init && terraform apply -auto-approve

quality-gates: fmt vet test build tf-validate
	echo "✅✅✅"

.PHONY: build