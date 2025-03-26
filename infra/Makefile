fmt:
	terraform fmt -recursive .

fmt-check:
	@echo "Checking Terraform formatting..."
	terraform fmt -recursive -check .

lint:
	cd base && terraform init && terraform validate
	cd env/dev && terraform init && terraform validate

apply-auto:
	@if [ -z "$(env)" ]; then echo "Error: env is not set. Please pass by name: env=<dev|staging|prod>."; exit 1; fi
	cd "env/${env}" && terraform init && terraform apply -auto-approve

apply:
	@if [ -z "$(env)" ]; then echo "Error: env is not set. Please pass by name: env=<dev|staging|prod>."; exit 1; fi
	cd "env/${env}" && terraform init && terraform apply

plan:
	@if [ -z "$(env)" ]; then echo "Error: env is not set. Please pass by name: env=<dev|staging|prod>."; exit 1; fi
	cd "env/${env}" && terraform init && terraform plan

quality-gates: fmt lint
	make plan env=dev
	echo "✅✅✅"

.PHONY: build