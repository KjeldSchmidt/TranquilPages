fmt:
	terraform fmt -recursive .

fmt-check:
	@echo "Checking Terraform formatting..."
	terraform fmt -recursive -check .

lint:
	cd base && terraform init && terraform validate
	cd env/dev && terraform init && terraform validate

tf-apply-auto:
	@if [ -z "$(env)" ]; then echo "Error: env is not set. Please pass by name: env=<dev|staging|prod>."; exit 1; fi
	cd "env/${env}" && terraform init && terraform apply -auto-approve

quality-gates: fmt lint
	echo "✅✅✅"

.PHONY: build