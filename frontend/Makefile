run:
	npm run ng -- serve --poll=2000

test-watch:
	npm run ng -- test --poll=2000

test:
	npm run ng -- test --no-watch --no-progress --browsers=ChromeHeadless

build:
	dotenvx run -f ../.env.local -- node scripts/generate-env.js
	npm run ng build

quality-gates: test

install-repo-dependencies:
	npm install

install-repo-dependencies-wsl:
	sudo apt-get update
	sudo apt-get install -y curl unzip xvfb libxi6 libgconf-2-4
	wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
	sudo apt install ./google-chrome-stable_current_amd64.deb