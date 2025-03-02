
FROM docker.arvancloud.ir/golang:1.24-alpine AS builder

WORKDIR /src

COPY . .

RUN go build -o /bin/publisher /src/publisher.go

FROM docker.arvancloud.ir/alpine:latest AS runtime

WORKDIR /app

COPY --from=builder /bin/publisher /app/publisher

CMD ["/app/publisher"]
