 .PHONY: run migrate build

run: migrate build
	docker compose up

build:
	docker-compose build --no-cache api