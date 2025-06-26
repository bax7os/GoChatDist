# Makefile para o projeto GoChatDist

# ==============================================================================
# VARIÁVEIS DE CONFIGURAÇÃO
# Usar variáveis torna o Makefile mais fácil de manter e adaptar.
# ==============================================================================

# Nomes dos binários de saída que serão criados na raiz do projeto.
CLIENT_BIN := client
SERVER_BIN := server

# Caminhos para os pacotes principais do cliente e do servidor.
CLIENT_DIR := ./cmd/client
SERVER_DIR := ./cmd/server

# Localização dos arquivos .proto e dos arquivos Go que serão gerados.
# Usar wildcards torna o processo automático se você adicionar mais arquivos .proto.
PROTO_DIR := proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
GENERATED_PB_FILES := $(patsubst %.proto,$(PROTO_DIR)/%.pb.go,$(notdir $(PROTO_FILES)))
GENERATED_GRPC_FILES := $(patsubst %.proto,$(PROTO_DIR)/%_grpc.pb.go,$(notdir $(PROTO_FILES)))

# ==============================================================================
# ALVOS (.PHONY)
# .PHONY evita que o Make confunda os alvos com nomes de arquivos reais.
# ==============================================================================

.PHONY: all proto tidy build build-client build-server run run-client run-server clean help

# O alvo 'help' exibe uma mensagem de ajuda. É uma boa prática para Makefiles.
help:
	@echo "Makefile para o projeto GoChatDist"
	@echo ""
	@echo "Uso: make [alvo]"
	@echo ""
	@echo "Alvos disponíveis:"
	@echo "  all           Constrói o cliente e o servidor (alvo padrão)."
	@echo "  proto         Gera o código Go a partir dos arquivos .proto."
	@echo "  tidy          Executa 'go mod tidy' para arrumar as dependências."
	@echo "  build         Constrói o cliente e o servidor (sinônimo de 'all')."
	@echo "  build-client  Constrói apenas o binário do cliente."
	@echo "  build-server  Constrói apenas o binário do servidor."
	@echo "  run           Executa o servidor e o cliente em paralelo para desenvolvimento."
	@echo "  run-server    Executa o servidor gRPC (porta 50051)."
	@echo "  run-client    Executa o cliente/servidor web (porta 8080)."
	@echo "  clean         Remove os binários e os arquivos gerados pelo Protobuf."
	@echo "  help          Exibe esta mensagem de ajuda."
	@echo ""

# ==============================================================================
# ALVOS DE CONSTRUÇÃO E EXECUÇÃO
# ==============================================================================

# O comando 'make' ou 'make all' irá construir ambos os binários.
# 'build' é um alias mais semântico para 'all'.
all: build
build: build-server build-client

# 1. Gera o código Go a partir dos arquivos .proto.
#    Este alvo só será executado se os arquivos .proto forem mais recentes que os gerados.
proto: $(GENERATED_PB_FILES) $(GENERATED_GRPC_FILES)
$(GENERATED_PB_FILES) $(GENERATED_GRPC_FILES): $(PROTO_FILES)
	@echo "--- Gerando código do Protobuf a partir de todos os arquivos em $(PROTO_DIR) ---"
	@protoc --go_out=. --go-grpc_out=. $(PROTO_FILES)

# 2. Arruma as dependências do módulo.
#    Depende de 'proto' para garantir que o código gerado seja incluído.
tidy: proto
	@echo "--- Arrumando as dependências do módulo (go mod tidy) ---"
	@go mod tidy

# 3. Constrói o binário do cliente.
#    Depende de 'tidy' para garantir que tudo esteja atualizado.
build-client: tidy
	@echo "--- Construindo binário do cliente: $(CLIENT_BIN) ---"
	cd $(CLIENT_DIR) && go build -o client

# 4. Constrói o binário do servidor gRPC.
#    Depende de 'tidy' para garantir que tudo esteja atualizado.
build-server: tidy
	@echo "--- Construindo binário do servidor gRPC: $(SERVER_BIN) ---"
	cd $(SERVER_DIR) && go build -o server

# Roda o servidor gRPC. Primeiro, garante que a versão mais recente está compilada.
run-server: build-server
	@echo "--- Executando servidor gRPC na porta 50051... (Pressione Ctrl+C para parar) ---"
	cd $(SERVER_DIR) && ./server

# Roda o cliente. Primeiro, garante que a versão mais recente está compilada.
run-client: build-client
	@echo "--- Executando cliente na porta 8080... (Pressione Ctrl+C para parar) ---"
	cd $(CLIENT_DIR) && ./client


# Roda ambos os serviços em paralelo para facilitar o desenvolvimento.
# O servidor é executado em segundo plano (&) e o cliente em primeiro plano.
# Ao fechar o cliente (Ctrl+C), o terminal também encerrará o servidor.
run: build
	@echo "--- Executando servidor em background e cliente em foreground ---"
	@./$(SERVER_BIN) & ./$(CLIENT_BIN)

# ==============================================================================
# ALVO DE LIMPEZA
# ==============================================================================

# Limpa os binários e os arquivos .go gerados pelo Protobuf.
# A correção aqui é garantir que os arquivos *_grpc.pb.go também sejam removidos.
clean:
	@echo "--- Limpando binários e arquivos gerados ---"
	@rm -f $(CLIENT_BIN) $(SERVER_BIN)
	@rm -f $(PROTO_DIR)/*.pb.go
	@rm -f $(PROTO_DIR)/*_grpc.pb.go
	@echo "--- Limpeza concluída ---"

