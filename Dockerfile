FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o smart-home .

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/smart-home .
COPY --from=builder /app/.env .

RUN chmod +x smart-home

ENV TZ=Europe/Athens
RUN ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD ["./smart-home"]
