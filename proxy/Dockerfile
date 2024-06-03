FROM golang:1.19.1-alpine3.15 AS builder 
WORKDIR /app
COPY ./ /app

RUN go mod tidy 
RUN go build -ldflags="-w -s" -o /app/server ./cmd/proxy

FROM alpine:3.15

COPY --from=builder /app/server /app/server 
COPY ./public /app/public
COPY ./.env /app/.env

WORKDIR /app 

ENTRYPOINT [ "/app/server" ]