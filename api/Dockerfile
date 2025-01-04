FROM golang:alpine AS builder 

RUN mkdir /build
COPY go.mod go.sum ./
RUN go mod download
ADD . /build

WORKDIR /build

RUN go build -o main .

#stage 2

FROM alpine:latest

RUN adduser -S -D -H -h /app appuser

USER appuser

COPY --from=builder /build/main /app/main

WORKDIR /app    

COPY .env /app/.env

EXPOSE 3000

CMD ["./main"]