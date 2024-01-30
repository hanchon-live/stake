.PHONY: stake

dev:
	@source /opt/homebrew/opt/nvm/nvm.sh && nvm use v20.10.0 && air

lint:
	@templ fmt .

test:
	@godotenv -f .env go test -v ./...

build:
	@./build.sh
