BINARY_NAME=Nourybot-Matrix.out

dev:
	go build -o ${BINARY_NAME} cmd/bot/*
	./${BINARY_NAME} --database "./db/nourybot.db"

build:
	go build -o ${BINARY_NAME} cmd/bot/*

rebuild:
	docker compose down
	docker-compose up --force-recreate --no-deps --build nourybot-matrix
