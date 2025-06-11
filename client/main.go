package main

import (
	"context"
	"gochatdist/messaging"
	pb "gochatdist/proto"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	username := "Carlos"

	// Inicia o subscriber para receber mensagens
	messaging.SubscribeToQueue(username)

	// Conecta ao gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Erro ao conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	// Envia uma mensagem
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.SendMessage(ctx, &pb.MessageRequest{
		Sender:   "Ana",
		Receiver: username,
		Content:  "Olá, Carlos!",
	})
	if err != nil {
		log.Fatalf("Erro ao enviar mensagem: %v", err)
	}

	log.Printf("Resposta do servidor: %s", resp.Status)

	// Mantém o programa rodando para ouvir mensagens
	select {}
}
