package server

import (
	"context"
	"fmt"
	"gochatdist/messaging"
	pb "gochatdist/proto"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
}

func NewChatServer() *ChatServer {
	return &ChatServer{}
}

func (s *ChatServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	fmt.Printf("Mensagem recebida de %s para %s: %s\n", req.Sender, req.Receiver, req.Content)

	// Publicar no RabbitMQ
	queueName := req.Receiver
	err := messaging.PublishMessage(fmt.Sprintf("%s: %s", req.Sender, req.Content), queueName)
	if err != nil {
		return nil, err
	}

	return &pb.MessageResponse{Status: "Mensagem publicada para " + queueName}, nil
}
