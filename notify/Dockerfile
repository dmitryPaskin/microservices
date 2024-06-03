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