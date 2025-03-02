FROM docker.arvancloud.ir/golang:1.24-alpine as builder

RUN apk add --no-cache git

WORKDIR /src

RUN go get github.com/julienschmidt/httprouter
RUN go get github.com/sirupsen/logrus
RUN go get github.com/streadway/amqp

COPY . .

RUN go build -o /bin/publisher ./publisher.go

FROM docker.arvanclout.ir/alpine:latest as runtime

COPY --from=builder /src/publisher /app/publisher

CMD ["/app/publisher"]
