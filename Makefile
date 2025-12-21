APP_NAME=pointservice

.PHONY: tidy build run up down

tidy:
	go mod tidy

build:
	docker compose build

up:build
	docker compose up

run: up

down:
	docker compose down -v