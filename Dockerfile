FROM golang:1.24.5 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build main.go

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go -o ./docs --parseDependency --parseInternal

FROM ubuntu:latest  
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y bash

WORKDIR /root/

COPY --from=builder /app/build .
COPY --from=builder /app/.env .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations

CMD ["./build"]