BINARY_NAME=Nourybot-Matrix.out

dev:
	go build -o ${BINARY_NAME} cmd/bot/*
	./${BINARY_NAME} --env="dev"

prod:
	go build -o ${BINARY_NAME} cmd/bot/*
	./${BINARY_NAME} --env="prod"

build:
	go build -o ${BINARY_NAME} cmd/bot/*

rebuild:
	docker compose down
	docker-compose up --force-recreate --no-deps --build nourybot-matrix
