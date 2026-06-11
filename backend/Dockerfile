FROM golang:1.22-alpine AS builder

WORKDIR /app

ENV GOSUMDB=off \
    GOPROXY=https://proxy.golang.org,direct

COPY go.mod ./
COPY *.go ./

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
ENTRYPOINT ["./server"]
