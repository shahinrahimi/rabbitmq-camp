package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type jsonResponse struct {
  Error bool
}


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

func main() {
  router := httprouter.New()

  router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Typ", "plain/text")    
    w.WriteHeader(http.StatusOK)
    _, err := w.Write([]byte("it's works!"))
    if err != nil {
      log.Println("Failed to write to Response")
    }
  })
  router.POST("/publish/:message", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    log.Println("Reciebed a post request at publish end")
    submit(w,r,p)
  })

  log.Println("Running publisher....")
  log.Fatal(http.ListenAndServe(":5000", router))
}

func submit(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

  message := p.ByName("message")
  log.Println("Recieved message: "+ message)

  url := fmt.Sprintf("amqp://%s:%s@%s:%s/",rabbit_user, rabbit_pass, rabbit_host, rabbit_port)
  log.Printf("Try connect to RabbitMQ with url: %s\n", url)

  conn, err := amqp.Dial(url)
  if err != nil {
    log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
  }
  defer conn.Close()
  
  ch, err := conn.Channel()

  if err != nil {
    log.Fatalf("%s: %s", "Failed to open a channel", err)
  }
  defer ch.Close()

  q, err := ch.QueueDeclare(
    "publisher_q",
    false, // durable
    false, // delete when used
    false, // exclusive
    false, // no-wait
    nil, // arguments
    )
  if err != nil {
    log.Fatalf("%s: %s", "Failed to declare a queue", err)
  }

  if err := ch.Publish(
    "", //exchange
    q.Name, // routing key
    false,  // mandatory
    false,  //immediate
    amqp.Publishing{
      ContentType: "text/plain",
      Body: []byte(message),
    }); err != nil {
      log.Fatalf("%s: %s", "Failed to publish a message", err)
  }

  log.Println("publish successfuly!")
}
