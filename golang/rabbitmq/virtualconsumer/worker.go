package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// for i := 0; i < 2; i++ {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	i := 2

	go func(i int) {
		for d := range msgs {
			// log.Printf("Received a message: %+v\n", d)
			// log.Printf("Received a message: %s", d.Body)
			if i == 0 {
				time.Sleep(2 * time.Second)
			}

			log.Println("Done worker", i)
			d.Ack(false)
		}
	}(i)
	// }

	forever := make(chan bool)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
