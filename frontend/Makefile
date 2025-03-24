run:
	ng serve --poll=2000

test-watch:
	ng test --poll=2000

test:
	ng test --watch=false

build:
	ng build

quality-gates: test

install-repo-dependencies-wsl:
	sudo apt-get update
	sudo apt-get install -y curl unzip xvfb libxi6 libgconf-2-4
	wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
	sudo apt install ./google-chrome-stable_current_amd64.deb