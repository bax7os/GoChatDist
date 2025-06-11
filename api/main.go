package main

import (
	"context"
	"encoding/json"
	pb "gochatdist/proto"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

type MessageRequest struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

type MessageResponse struct {
	Status string `json:"status"`
}

var grpcClient pb.ChatServiceClient

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// CORS e Content-Type
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Chama gRPC SendMessage
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcClient.SendMessage(ctx, &pb.MessageRequest{
		Sender:   req.Sender,
		Receiver: req.Receiver,
		Content:  req.Content,
	})
	if err != nil {
		http.Error(w, "Erro ao enviar mensagem via gRPC: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(MessageResponse{Status: resp.Status})
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Erro ao conectar gRPC: %v", err)
	}
	defer conn.Close()

	grpcClient = pb.NewChatServiceClient(conn)

	http.HandleFunc("/sendMessage", sendMessageHandler)

	log.Println("Servidor HTTP escutando na porta 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
