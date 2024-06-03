# FROM golang:1.19.1 AS builder 

# WORKDIR /app
# COPY ./ /app

# RUN apt-get update && apt-get install -y gcc libc-dev librdkafka-dev

# RUN go mod tidy 
# RUN CGO_ENABLED=1 go build -ldflags="-w -s" -o /app/server ./cmd/app

# FROM alpine:3.15

# RUN apk update && apk add --no-cache librdkafka

# COPY --from=builder /app/server /app/server 
# COPY ./.env /app/.env

# WORKDIR /app 

# ENTRYPOINT [ "/app/server" ]

# FROM golang:1.19.1 AS builder 

# WORKDIR /app
# COPY ./ /app

# RUN apt-get update && apt-get install -y gcc libc-dev librdkafka-dev
# RUN go mod tidy 
# RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/app

# FROM debian:bullseye-slim

# RUN apt-get update && apt-get install -y librdkafka-dev && rm -rf /var/lib/apt/lists/*

# COPY --from=builder /app/server /app/server
# COPY ./.env /app/.env

# WORKDIR /app

# ENTRYPOINT [ "/app/server" ]

# FROM golang:alpine AS builder

# RUN apk update && apk add --no-cache git gcc libc-dev gcompat

# WORKDIR /app
# COPY ./ /app

# ENV CGO_ENABLED=1
# RUN go build -o /app/server ./cmd/app

# FROM alpine

# RUN apk update && apk add --no-cache librdkafka gcompat

# COPY --from=builder /app/server /app/server 
# COPY ./.env /app/.env

# WORKDIR /app 

# ENTRYPOINT [ "/app/server" ]

# IT WORKS
FROM golang:1.19-buster AS builder 

WORKDIR /app 
COPY ./ /app

RUN apt-get update && apt-get install -y gcc libc-dev librdkafka-dev 
RUN go mod tidy 
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/app

FROM debian:buster

RUN apt-get update && \
apt-get install -y ca-certificates librdkafka-dev && \
update-ca-certificates && \
rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server 
COPY ./.env /app/.env

WORKDIR /app

ENTRYPOINT ["/app/server"]


