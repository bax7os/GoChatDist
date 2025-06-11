package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Message struct {
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func SaveMessage(username string, msg Message) error {
	filePath := fmt.Sprintf("storage/%s.json", username)

	// Lê mensagens existentes
	var messages []Message
	data, err := os.ReadFile(filePath)
	if err == nil {
		json.Unmarshal(data, &messages)
	}

	// Adiciona a nova mensagem
	messages = append(messages, msg)

	// Escreve de volta no arquivo
	newData, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, newData, 0644)
}
