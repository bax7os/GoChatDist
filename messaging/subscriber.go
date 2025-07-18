// Em messaging/subscriber.go
package messaging

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func SubscribeToQueue(queueName string, handler func(msgBody []byte)) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Printf("Erro ao conectar ao RabbitMQ: %v", err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Erro ao abrir canal: %v", err)
		return
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable: Garante que o subscriber se conecte à mesma fila durável.
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		log.Printf("Erro ao declarar fila: %v", err)
		return
	}

	
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Erro ao registrar consumidor: %v", err)
		return
	}

	fmt.Printf("Subscriber iniciado para a fila durável: %s\n", queueName)

	for d := range msgs {
		handler(d.Body)
	}

	log.Printf("Subscriber para %s foi encerrado.", queueName)
}
