env ?= local

build:
	cd backend && make build env=${env}
	cd frontend && make build env=${env}

test:
	cd backend && make test
	cd frontend && make test

fmt:
	cd backend && make fmt
	cd infra && make fmt

fmt-check:
	cd backend && make fmt-check
	cd infra && make fmt-check

run:
	cd frontend && make run &
	cd backend && make run &

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

create-env-file:
	@if [ -z "$(env)" ]; then echo "Error: env is not set. Please pass by name: env=<dev|staging|prod>."; exit 1; fi
	rm ".env.${env}"
	touch ".env.${env}"
	printf "FRONTEND_URL='" >> ".env.${env}"
	(cd "infra/env/${env}" && terraform output -raw frontend_url) >> ".env.${env}"
	echo "'" >> ".env.${env}"

	printf "BACKEND_URL='" >> ".env.${env}"
	(cd "infra/env/${env}" && terraform output -raw backend_url) >> ".env.${env}"
	echo "'" >> ".env.${env}"

reset-local-db:
	sudo docker stop mongodb
	sudo docker rm mongodb
	sudo docker run -d --name mongodb -p 27017:27017 mongo:latest

install-repo-dependencies:
	git config core.hooksPath ./.githooks
	curl -fsS https://dotenvx.sh | sudo sh
	cd backend && make install-repo-dependencies
	cd frontend && make install-repo-dependencies

install-repo-dependencies-ci:
	curl -fsS https://dotenvx.sh | sudo sh
	cd backend && make install-repo-dependencies
	cd frontend && make install-repo-dependencies

.PHONY: build