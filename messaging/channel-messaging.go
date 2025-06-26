// Crie um novo arquivo ou adicione a um existente no pacote messaging
package messaging

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// getChannel se conecta ao RabbitMQ e retorna um canal
func getChannel() (*amqp.Channel, func()) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return nil, func() {}
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		conn.Close()
		return nil, func() {}
	}

	// Retorna o canal e uma função para fechar tudo
	closeFunc := func() {
		ch.Close()
		conn.Close()
	}

	return ch, closeFunc
}

// PublishToChannel publica uma mensagem em uma exchange do tipo fanout.
func PublishToChannel(channelName string, message string) error {
	ch, close := getChannel()
	if ch == nil {
		return fmt.Errorf("could not get rabbitmq channel")
	}
	defer close()

	exchangeName := "exchange." + channelName

	// Declara a exchange do tipo fanout. Se já existir, não faz nada.
	err := ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// Publica a mensagem na exchange
	err = ch.Publish(
		exchangeName, // exchange
		"",           // routing key (não é usada em fanout)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf(" [x] Sent to channel %s: %s", channelName, message)
	return nil
}

// BindToChannel liga (bind) a fila de um usuário a uma exchange de canal.
func BindToChannel(queueName, channelName string) error {
	ch, close := getChannel()
	if ch == nil {
		return fmt.Errorf("could not get rabbitmq channel")
	}
	defer close()

	exchangeName := "exchange." + channelName

	// Garante que a exchange existe
	_ = ch.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil)

	// Liga a fila à exchange
	err := ch.QueueBind(
		queueName,    // queue name
		"",           // routing key (não usada em fanout)
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind a queue: %w", err)
	}
	
	log.Printf("Queue %s bound to exchange %s", queueName, exchangeName)
	return nil
}

// UnbindFromChannel desfaz a ligação da fila com a exchange.
func UnbindFromChannel(queueName, channelName string) error {
	ch, close := getChannel()
	if ch == nil {
		return fmt.Errorf("could not get rabbitmq channel")
	}
	defer close()

	exchangeName := "exchange." + channelName

	err := ch.QueueUnbind(
		queueName,    // queue name
		"",           // routing key
		exchangeName, // exchange
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to unbind queue: %w", err)
	}
	
	log.Printf("Queue %s unbound from exchange %s", queueName, exchangeName)
	return nil
}
