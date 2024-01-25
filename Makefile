.PHONY: stake

dev:
	@source /opt/homebrew/opt/nvm/nvm.sh && nvm use v20.10.0 && air

lint:
	@templ fmt .

build:
	@./build.sh
