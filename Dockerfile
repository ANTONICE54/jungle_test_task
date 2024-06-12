#Build stage
FROM golang:1.21.4-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/main.go

#Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
RUN mkdir /app/images




EXPOSE 8080
