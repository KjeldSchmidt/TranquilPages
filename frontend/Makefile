env ?= local

run:
	npm run ng -- serve --poll=2000

test-watch:
	npm run ng -- test --poll=2000

test:
	npm run ng -- test --no-watch --no-progress --browsers=ChromeHeadless

build:
	@if [ "$(CI)" = "true" ] && [ "$(env)" = "local" ]; \
 		then echo ".env.local cannot be used to build in a pipeline environment! Please pass env by name: env=<dev|staging|prod>."; \
 		exit 1; \
 	fi
	dotenvx run -f ../.env.${env} -- node scripts/generate-env.js
	npm run ng build

quality-gates: test

deploy:
	@if [ "$(env)" = "local" ]; then echo "Error: env is not set. Please pass by name: env=<dev|staging|prod>."; exit 1; fi

	@export $(shell cat ../.env.${env}) && az storage blob upload-batch \
		--auth-mode key \
		--account-name $$FRONTEND_STORAGE_ACCOUNT \
		--destination '$$web' \
		--source dist/tranquil-pages/browser \
		--overwrite

	@export $(shell cat ../.env.${env}) && az storage blob delete-batch \
		--auth-mode key \
		--source '$$web' \
		--account-name $$FRONTEND_STORAGE_ACCOUNT \
		--if-unmodified-since $(shell date -u --date='5 minutes ago' +%Y-%m-%dT%H:%MZ)


install-repo-dependencies:
	npm install

install-repo-dependencies-wsl:
	sudo apt-get update
	sudo apt-get install -y curl unzip xvfb libxi6 libgconf-2-4
	wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
	sudo apt install ./google-chrome-stable_current_amd64.deb