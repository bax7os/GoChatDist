package messaging

import (
	"github.com/streadway/amqp"
)

//http://localhost:15672/#/queues


func PublishMessage(message string, queueName string) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	
	err = ch.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return err
}
