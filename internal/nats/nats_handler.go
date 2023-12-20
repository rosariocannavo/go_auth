package nats

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/rosariocannavo/go_auth/config"
)

// startup a nats connection and let it be accessible by NatsConnection
type NatsConnectionHandler struct {
	*nats.Conn
}

var NatsConnection *NatsConnectionHandler

func init() {
	nc, err := nats.Connect(config.NatsURL)
	if err != nil {
		fmt.Println("Nats: impossible to connect")
		log.Fatal(err)

	}

	NatsConnection = &NatsConnectionHandler{nc}

}

func (n *NatsConnectionHandler) close() {
	if n.Conn != nil {
		n.Conn.Close()
	}
}

func (n *NatsConnectionHandler) PublishMessage(message string) {
	err := n.Publish(config.NatsSubject, []byte(message))
	if err != nil {
		fmt.Println("Nats: impossible to publish message")
		log.Fatal(err)
		n.close()
	} else {
		log.Print("Nats published")
	}
}
