package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/assembla/cony"
)

func MustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s missing", key)
	}
	return v
}

func passMsg(data amqp.Delivery) {
	log.Info(data)
}

func main() {

	// Construct new client with the flag url
	// and default backoff policy
	cli := cony.NewClient(
		cony.URL("amqp://user:bitnami@rabbitmq:5672/"),
		cony.Backoff(cony.DefaultBackoff),
	)

	// Declarations
	// The queue name will be supplied by the AMQP server
	que := &cony.Queue{
		Name:       MustGetEnv("Q_NAME"),
		AutoDelete: false,
	}

	exc := cony.Exchange{
		Name:       MustGetEnv("EXCHANGE_NAME"),
		Kind:       "direct",
		Durable:    true,
		AutoDelete: false,
	}
	bnd := cony.Binding{
		Queue:    que,
		Exchange: exc,
		Key:      MustGetEnv("ROUTING_KEY"),
	}
	cli.Declare([]cony.Declaration{
		cony.DeclareQueue(que),
		cony.DeclareExchange(exc),
		cony.DeclareBinding(bnd),
	})

	// Declare and register a consumer
	cns := cony.NewConsumer(
		que,
		cony.Qos(1),
	)
	cli.Consume(cns)
	for cli.Loop() {
		select {
		case msg := <-cns.Deliveries():
			go passMsg(msg)
			// If when we built the consumer we didn't use
			// the "cony.AutoAck()" option this is where we'd
			// have to call the "amqp.Deliveries" methods "Ack",
			// "Nack", "Reject"
			//
			// msg.Ack(false)
			// msg.Nack(false)
			// msg.Reject(false)
			msg.Ack(false)
		case err := <-cns.Errors():
			log.Infof("Consumer error: %v\n", err)
		case err := <-cli.Errors():
			log.Infof("Client error: %v\n", err)
		}
	}
}
