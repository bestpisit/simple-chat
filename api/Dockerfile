# Use the official Golang image to create a build artifact.
FROM golang:1.21 as builder

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Copy the go.mod and go.sum file and download the dependencies.
COPY go.* ./
RUN go mod download

# Copy the rest of the application's code.
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:3
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /server

CMD ["/server"]
