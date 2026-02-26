FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/finalyandex ./main.go

FROM alpine:3.22
WORKDIR /app

COPY --from=builder /out/finalyandex /app/finalyandex
COPY web /app/web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

EXPOSE 7540

CMD ["/app/finalyandex"]
