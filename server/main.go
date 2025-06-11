package main

import (
	"fmt"
	pb "gochatdist/proto"
	server "gochatdist/server/chat"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, server.NewChatServer())

	fmt.Println("Servidor gRPC escutando na porta 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Erro ao servir: %v", err)
	}
}
