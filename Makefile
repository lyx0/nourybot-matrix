rebuild:
	docker compose down
	docker-compose up --force-recreate --no-deps --build nourybot-matrix