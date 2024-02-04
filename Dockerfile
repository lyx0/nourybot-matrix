# Start from golang base image
FROM golang:alpine3.19

RUN apk add --no-cache git ca-certificates build-base su-exec olm-dev sqlite

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .
#COPY .env .
COPY ./db/nourybot.db .

# Download all the dependencies
RUN go get -d -v ./...
Run go build -o "Nourybot.out" ./cmd/bot

# Run the executable
CMD [ "./Nourybot.out" ]

