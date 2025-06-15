package server

import (
	"context"
	"encoding/json"
	"fmt"
	"gochatdist/messaging"
	pb "gochatdist/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
}

func NewChatServer() *ChatServer {
	return &ChatServer{}
}

// formata a mensagem e a envia para o rabbitmq passando o sender como username 
func (s *ChatServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
    // Validação
    if req.Sender == "" || req.Receiver == "" {
        return nil, status.Errorf(codes.InvalidArgument, "sender and receiver cannot be empty")
    }

    // formatar como JSON para permanencia no rabbitmq
    msg := struct {
        Sender  string `json:"sender"`
        Content string `json:"content"`
    }{
        Sender:  req.Sender,
        Content: req.Content,
    }

    msgBytes, err := json.Marshal(msg)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to marshal message: %v", err)
    }

    // Publicar na fila
    if err := messaging.PublishMessage(string(msgBytes), req.Receiver); err != nil {
        return nil, status.Errorf(codes.Internal, "failed to publish message: %v", err)
    }

    return &pb.MessageResponse{
        Status: fmt.Sprintf("Message published to %s", req.Receiver),
    }, nil
}