package broker

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQBroker(url string) (*RabbitMQBroker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQBroker{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQBroker) Publish(queue string, message string) error {
	_, err := r.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (r *RabbitMQBroker) Subscribe(queue string, handler func(message string) error) error {
	_, err := r.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err := handler(string(d.Body)); err != nil {
				log.Printf("error handling message: %v", err)
			}
		}
	}()

	return nil
}

func (r *RabbitMQBroker) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}

	if r.conn != nil {
		r.conn.Close()
	}
	return nil
}
