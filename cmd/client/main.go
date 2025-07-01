package main

import (
	"context"
	"fmt"
	"gochatdist/messaging"
	pb "gochatdist/proto"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients   = make(map[string]*websocket.Conn)
	clientsMu sync.Mutex
) // variável do cliente

func handleConnections(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username é obrigatório", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Erro no upgrade: %v", err)
		return
	}
	defer ws.Close()

	clientsMu.Lock()
	clients[username] = ws
	clientsMu.Unlock()
	log.Printf("Cliente conectado: %s", username)

	wsCallback := func(msgBody []byte) {
		clientsMu.Lock()
		defer clientsMu.Unlock()
		if conn, ok := clients[username]; ok {
			if err := conn.WriteMessage(websocket.TextMessage, msgBody); err != nil {
				log.Printf("Erro ao enviar msg para o websocket de %s: %v", username, err)
			}
		}
	}
	
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Erro ao conectar ao gRPC: %v", err)
		return
	}
	defer conn.Close()
	grpcClient := pb.NewChatServiceClient(conn)

	// --- PASSO 1: REGISTRAR O USUÁRIO NO SERVIÇO DE NOMES ---
	ctxReg, cancelReg := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = grpcClient.RegisterUser(ctxReg, &pb.RegisterUserRequest{Username: username})
	if err != nil {
		log.Printf("Falha ao registrar usuário %s: %v", username, err)
		wsCallback([]byte("Erro: não foi possível registrar no servidor."))
		return
	}
	cancelReg()
	// --------------------------------------------------------

	go messaging.SubscribeToQueue(username, wsCallback)
	messaging.BindToChannel(username, "geral")
	wsCallback([]byte("Você entrou no canal #geral. Para ver os comandos do servidor digite @Servidor comandos. Use /join #nome para entrar em um chat, /leave #nome para sair de um chat, ou #nome <msg> para enviar uma mensagem para um chat específico..."))


	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Cliente %s desconectado: %v", username, err)
			clientsMu.Lock()
			delete(clients, username)
			clientsMu.Unlock()
			break
		}

		input := string(message)
		parts := strings.SplitN(input, " ", 2)
		command := parts[0]

		if strings.HasPrefix(command, "@") {
			if len(parts) < 2 { continue }
			receiver := strings.TrimPrefix(command, "@")
			content := parts[1]
			complete_message := ""
			// --- PASSO 2: TRATAR O ERRO DE USUÁRIO INEXISTENTE E MENSAGEM PARA O SERVIDOR ---
			if receiver == "Servidor" {
				complete_message = content // content = "lista de usuários", "comandos"
			} else {
				complete_message = fmt.Sprintf("(DM de %s) %s", username, content)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := grpcClient.SendMessage(ctx, &pb.MessageRequest{Sender: username, Receiver: receiver, Content: complete_message})
			if err != nil {
				log.Printf("Erro ao enviar via gRPC: %v", err)
				// Envia o feedback de erro para o usuário!
				wsCallback([]byte(fmt.Sprintf("Sistema: %v", err)))
			}
			cancel()
			// ----------------------------------------------------
		// Envio de mensagem para grupo
		} else if strings.HasPrefix(command, "#") {
			// (código dos canais permanece igual)
			if len(parts) < 2 { continue }
			channelName := strings.TrimPrefix(command, "#")
			content := parts[1]
			fullMessage := fmt.Sprintf("[#%s] %s: %s", channelName, username, content)
			err := messaging.PublishToChannel(channelName, fullMessage)
			if err != nil { log.Printf("Erro ao publicar no canal %s: %v", channelName, err) }
		// Entrar e sair de grupos
		} else if strings.HasPrefix(command, "/") {
			// (código de /join e /leave permanece igual)
			if len(parts) < 2 { continue }
			channelName := strings.TrimPrefix(parts[1], "#")
			switch command {
			case "/join":
				messaging.BindToChannel(username, channelName)
				wsCallback([]byte(fmt.Sprintf("Inscrito em #%s", channelName)))
			case "/leave":
				messaging.UnbindFromChannel(username, channelName)
				wsCallback([]byte(fmt.Sprintf("Saiu de #%s", channelName)))
			}
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleConnections)
	log.Println("Servidor Web escutando na porta :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor web: %v", err)
	}
}
