package messaging

import (
	"encoding/json"
	"fmt"
	"gochatdist/storage"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func SubscribeToQueue(queueName string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir canal: %v", err)
	}

	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Erro ao declarar fila: %v", err)
	}

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Erro ao registrar consumidor: %v", err)
	}

	go func() {
		for d := range msgs {
			fmt.Printf("Nova mensagem: %s\n", d.Body)

			var content struct {
				Sender  string `json:"sender"`
				Content string `json:"content"`
			}
			// Aqui para facilitar, suponho que a mensagem esteja no formato JSON
			// Se estiver como texto simples, pode simplificar
			err := json.Unmarshal(d.Body, &content)
			if err != nil {
				// Se for texto plano, pode ignorar o erro e salvar como conteúdo bruto
				content.Content = string(d.Body)
			}

			// Salvar mensagem recebida
			storage.SaveMessage(queueName, storage.Message{
				Sender:    content.Sender,
				Receiver:  queueName,
				Content:   content.Content,
				Timestamp: time.Now(),
			})
		}
	}()

	fmt.Printf("Escutando mensagens na fila: %s\n", queueName)
}
