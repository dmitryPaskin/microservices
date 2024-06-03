FROM golang:1.19.1-alpine3.15 AS builder 
WORKDIR /app
COPY ./ /app

RUN go mod tidy 
RUN go build -ldflags="-w -s" -o /app/server ./cmd/app

FROM alpine:3.15

COPY --from=builder /app/server /app/server 
COPY ./.env /app/.env

WORKDIR /app 

ENTRYPOINT [ "/app/server" ]