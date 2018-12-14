FROM golang:1.11 AS builder

RUN mkdir /app
WORKDIR /app

# Download (and cache) dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# Create image
FROM scratch
COPY --from=builder /app/main ./
ENTRYPOINT ["./main"]
