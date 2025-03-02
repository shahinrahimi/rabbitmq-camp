package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
  rabbitHost = "RABBIT_HOST"
  rabbitPort = "RABBIT_PORT"
  rabbitUser = "RABBIT_USER"
  rabbitPass = "RABBIT_PASS"
)

var rabbit_host = os.Getenv(rabbitHost)
var rabbit_port = os.Getenv(rabbitPort)
var rabbit_user = os.Getenv(rabbitUser)
var rabbit_pass = os.Getenv(rabbitPass)

func main(){
 consume() 
}

func consume() {
  url := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbit_user, rabbit_pass, rabbit_host, rabbit_port)

  conn, err := amqp.Dial(url)
  if err != nil {
    log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
  }
  defer conn.Close()

  ch, err := conn.Channel()
  if err != nil {
    log.Fatalf("%s:%s", "Failed to open a channel", err)
  }
  defer ch.Close()

  q, err := ch.QueueDeclare(
    "publisher_queue",
    false, // durable
    false, // delete when used
    false, // exlusive
    false, // no wait
    nil,   // arguments
    )
  if err != nil {
    log.Fatalf("%s:%s", "Failed to declare a queue", err)
  }

  log.Println("Channel and Queue established")

  msgs, err := ch.Consume(
    q.Name,
    "",    // consumer
    false, // auto-ack
    false, // exlusive
    false, // no-local
    false, // no-wait
    nil,
    )

  if err != nil {
    log.Fatalf("%s:%s", "Failed to register consumer", err)
  }
  
  forever := make(chan bool)
  go func() {
    for d := range msgs {
      log.Printf("Received a message: %s", d.Body)
      d.Ack(false)
    }
  }()

  log.Println("Running...")
  <-forever


}
