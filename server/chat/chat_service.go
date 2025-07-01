package server

import (
	"context"
	"fmt"
	"gochatdist/messaging"
	pb "gochatdist/proto"
	"strings"
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
	comands := "lista de usuários, comandos"


	s.mu.RLock()
	_, exists := s.registeredUsers[receiver]
	s.mu.RUnlock()

	// Caso o destinatário seja o servidor:
	if receiver == "Servidor" {
		usuario := sender
		receiver = usuario
		sender = "Servidor"
		switch content { 
			case "lista de usuários":
				fmt.Printf("Enviando lista de usuário existentes para %s\n", usuario)
				var builder strings.Builder
				builder.WriteString("Lista de usuário existentes: ")
				for key,_ := range s.registeredUsers {
					builder.WriteString(key)
					builder.WriteString(", ")
				}
				content = strings.TrimSuffix(builder.String(), ", ")
			case "comandos":
				fmt.Printf("Enviando comandos do servidor para: %s\n", usuario)
				content = fmt.Sprintf("Comandos do servidor: %s\n", comands)
			default:
				fmt.Printf("Comando não reconhecido: %s\n", content)
				
				return nil, fmt.Errorf("comando '%s' não reconhecido", content)
		}
	} else if !exists {
		fmt.Printf("Tentativa de envio para usuário inexistente: %s\n", receiver)
	
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
