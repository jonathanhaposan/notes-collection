package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/febytanzil/gobroker"
	"github.com/febytanzil/gobroker/pubsub"
)

func main() {
	p := pubsub.NewPublisher(gobroker.RabbitMQ, pubsub.RabbitMQAMPQ("amqp://admin:admin@localhost:5672/", "/"))

	// ticker := time.NewTicker(1 * time.Millisecond)
	go func() {
		// for t := range ticker.C {
		// 	p.Publish("test.fanout", "msg"+t.String())
		// }
		for i := 0; i < 500; i++ {
			p.Publish("test.fanout", "msg"+time.Now().String())
		}
	}()

	s := pubsub.NewSubscriber(gobroker.RabbitMQ, []*pubsub.SubHandler{
		{
			Name:       "test",
			Topic:      "test.fanout",
			Handler:    testRMQ,
			MaxRequeue: 10,
			Concurrent: 2,
		},
	}, pubsub.RabbitMQAMPQ("amqp://admin:admin@localhost:5672/", "/"))

	s.Start()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	s.Stop()
	log.Println("shutting down")
	os.Exit(0)
}

func testRMQ(msg *gobroker.Message) error {
	log.Println("consume rabbitmq", string(msg.Body))
	return nil
}
