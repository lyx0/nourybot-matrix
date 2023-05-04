BINARY_NAME=NourybotMatrix.out

cup:
	sudo docker compose up

xd:
	go build -o ./bin/${BINARY_NAME} && ./bin/${BINARY_NAME} --homeserver="matrix.xxx" --username="xxx" --password="xxx"
