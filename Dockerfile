FROM golang:latest AS builder
COPY . /src

WORKDIR /src

RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static -s -w" -o compressboom

FROM alpine:latest

COPY --from=builder /src/compressboom /app/compressboom

RUN chmod +x /app/compressboom

WORKDIR /app

CMD ["/app/compressboom"]
