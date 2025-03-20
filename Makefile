build:
	cd backend && make build

test:
	cd backend && make test

fmt:
	cd backend && make fmt
	cd infra && make fmt

fmt-check:
	cd backend && make fmt-check
	cd infra && make fmt-check

run:
	cd backend && make run

lint:
	cd backend && make lint
	cd infra && make lint

tf-apply-auto:
	@if [ -z "$(env)" ]; then echo "Error: env is not set. Please pass by name: env=<dev|staging|prod>."; exit 1; fi
	cd "infra/" && make apply-auto "env=${env}"

quality-gates: fmt lint test build
	echo "✅✅✅"

build-image:
	cd backend && make build-image

push-image:
	cd backend && make push-image

reset-local-db:
	sudo docker stop mongodb
	sudo docker rm mongodb
	sudo docker run -d --name mongodb -p 27017:27017 mongo:latest

install-repo-dependencies:
	git config core.hooksPath ./.githooks
	curl -fsS https://dotenvx.sh | sudo sh

.PHONY: build