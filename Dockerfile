# Start from golang base image
FROM golang:alpine3.19

RUN apk add --no-cache git ca-certificates build-base su-exec olm-dev sqlite


# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...
RUN go get -u maunium.net/go/mautrix

# Build the Go app
RUN go build .

# Run the executable
CMD [ "./nourybot-matrix" ]

