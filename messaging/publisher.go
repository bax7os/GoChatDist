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

	// Declara a fila, garantindo que os parâmetros são os mesmos do subscriber.
	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable: true -> AQUI ESTÁ A CORREÇÃO. Agora é consistente.
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// Publica a mensagem como persistente.
	err = ch.Publish(
		"",        // exchange
		queueName, // routing key (o nome da fila)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Garante que a mensagem sobreviva a reinicializações.
			ContentType:  "text/plain",
			Body:         []byte(message),
		},
	)
	return err
}