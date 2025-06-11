# Variáveis
PROTO_DIR=proto
PROTO_FILE=$(PROTO_DIR)/chat.proto
GO_OUT=.
SERVER_DIR=server
CLIENT_DIR=client
API_DIR=api

# Comando para gerar código Go do proto
proto:
	protoc --go_out=$(GO_OUT) --go-grpc_out=$(GO_OUT) $(PROTO_FILE)

# Build do servidor gRPC
build-server:
	cd $(SERVER_DIR) && go build -o server

# Build do cliente gRPC
build-client:
	cd $(CLIENT_DIR) && go build -o client

# Build da API HTTP
build-api:
	cd $(API_DIR) && go build -o api

# Executar servidor gRPC
run-server:
	cd $(SERVER_DIR) && ./server

# Executar cliente gRPC
run-client:
	cd $(CLIENT_DIR) && ./client

# Executar API HTTP
run-api:
	cd $(API_DIR) && ./api

# Tudo junto: gerar código + build tudo
build-all: proto build-server build-client build-api

.PHONY: proto build-server build-client build-api run-server run-client run-api build-all
