package server

import (
	"context"
	"fmt"
	"gochatdist/messaging"
	pb "gochatdist/proto"
	"sync"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	// Este map vai funcionar como nosso "banco de dados" de usuários.
	registeredUsers map[string]bool
	mu              sync.RWMutex // Usamos RWMutex para permitir múltiplas leituras simultâneas.
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		registeredUsers: make(map[string]bool),
	}
}

// RegisterUser adiciona um novo usuário ao nosso registro.
func (s *ChatServer) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	username := req.GetUsername()
	s.mu.Lock()
	s.registeredUsers[username] = true
	s.mu.Unlock()
	fmt.Printf("Usuário registrado: %s\n", username)
	return &pb.RegisterUserResponse{Status: "Usuário " + username + " registrado com sucesso."}, nil
}

// SendMessage agora verifica se o destinatário existe antes de enviar.
func (s *ChatServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	sender := req.GetSender()
	receiver := req.GetReceiver()
	content := req.GetContent()

	// --- VERIFICAÇÃO DE EXISTÊNCIA ---
	s.mu.RLock()
	_, exists := s.registeredUsers[receiver]
	s.mu.RUnlock()

	if !exists {
		fmt.Printf("Tentativa de envio para usuário inexistente: %s\n", receiver)
		// Retorna um erro que o cliente pode tratar.
		return nil, fmt.Errorf("usuário '%s' não existe", receiver)
	}
	// ---------------------------------

	fmt.Printf("Mensagem de %s para %s: %s\n", sender, receiver, content)

	err := messaging.PublishMessage(content, receiver)
	if err != nil {
		return nil, fmt.Errorf("falha ao encaminhar mensagem para %s: %v", receiver, err)
	}

	return &pb.MessageResponse{Status: "Mensagem encaminhada para a fila de " + receiver}, nil
}
