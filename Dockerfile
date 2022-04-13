FROM golang:1.17 AS builder
WORKDIR /go/src
COPY . .
RUN cp .env.example .env
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o WALLET ./cmd/server

FROM alpine:latest
RUN apk add --no-cache bash coreutils grep sed ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/.env .
COPY --from=builder /go/src/WALLET .
CMD ["./WALLET"]

EXPOSE 80