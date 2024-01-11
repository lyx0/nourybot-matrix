BINARY_NAME=Nourybot-Matrix.out

xd:
	go build -o ${BINARY_NAME} cmd/bot/*
	./${BINARY_NAME} --env="dev"

dev:
	./${BINARY_NAME} --env="dev"

prod:
	./${BINARY_NAME} --env="prod"

build:
	go build -o ${BINARY_NAME} cmd/bot/*

rebuild:
	docker compose down
	docker-compose up --force-recreate --no-deps --build nourybot-matrix
