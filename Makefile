BINARY_NAME=NourybotMatrix.out

cup:
	sudo docker compose up

xd:
	cd cmd/nourybot && go build -o ${BINARY_NAME} && ./${BINARY_NAME}