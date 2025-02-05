FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o smart-home .

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/smart-home .
COPY --from=builder /app/config.json .

RUN chmod +x smart-home

CMD ["./smart-home"]
