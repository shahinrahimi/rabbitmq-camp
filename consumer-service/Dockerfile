FROM docker.arvancloud.ir/golang:1.24-alpine AS builder

WORKDIR /src

COPY . .

RUN go build -o /bin/consumer /src/consumer.go

FROM docker.arvancloud.ir/alpine:latest AS runtime

WORKDIR /app

COPY --from=builder /bin/consumer /app/consumer

CMD ["/app/consumer"]
