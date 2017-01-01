package rabbitmq

import (
	"cqrs"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

//Client rabbitmq
type Client struct {
	conn *amqp.Connection
}

//NewClient returns a Client to acces to rabbitmq
func NewClient(username, password, host string, port int) (*Client, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port))
	return &Client{
		conn: conn,
	}, err
}

//Publish a event
func (c *Client) Publish(event cqrs.Event, bucket, subset string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(
		bucket,   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return err
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = ch.Publish(
		bucket, // exchange
		subset, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		},
	)

	return err
}
