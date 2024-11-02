# FROM golang:1.21-alpine
# FROM golang:1.23
FROM golang:1.23-alpine

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Initialize go.mod and download dependencies
RUN go mod download && \
    go mod tidy

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

EXPOSE 3000

CMD ["./app"]
