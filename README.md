Programa desenvolvido e testando sobre Linux (ubuntu/arch). Todos os passos a seguir são pensados para serem executados nesse ambiente.<br>
Instale as bibliotecas do projeto:<br>
#### Golang
No link https://go.dev/doc/install faço o download do .tar
No terminal use os seguintes comandos:<br>
`rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.4.linux-amd64.tar.gz`<br>
`export PATH=$PATH:/usr/local/go/bin`<br>

Para verificar se o Golang foi instalado com sucesso utilize:<br>
`go version`

#### gRPC
Primeiro instale o compilador protoc:<br>
`sudo apt install -y protobuf-compiler`<br>
`protoc --version` A versão deve ser acima da 3.0, caso não seja siga os passos no link https://protobuf.dev/installation/<br>

Então instale o gRPC usando os seguintes comandos:<br>
`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`<br>
`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

Caso queira, para testar a instalação siga os passo no link https://grpc.io/docs/languages/go/quickstart/

#### RabbitMQ
Para instalar o RabbitMQ use o script abaixo. Para mais informações visite o link: https://www.rabbitmq.com/docs/install-debian#apt-quick-start-cloudsmith

```bash
#!/bin/sh

sudo apt-get install curl gnupg apt-transport-https -y

## Team RabbitMQ's main signing key
curl -1sLf "https://keys.openpgp.org/vks/v1/by-fingerprint/0A9AF2115F4687BD29803A206B73A36E6026DFCA" | sudo gpg --dearmor | sudo tee /usr/share/keyrings/com.rabbitmq.team.gpg > /dev/null
## Community mirror of Cloudsmith: modern Erlang repository
curl -1sLf https://github.com/rabbitmq/signing-keys/releases/download/3.0/cloudsmith.rabbitmq-erlang.E495BB49CC4BBE5B.key | sudo gpg --dearmor | sudo tee /usr/share/keyrings/rabbitmq.E495BB49CC4BBE5B.gpg > /dev/null
## Community mirror of Cloudsmith: RabbitMQ repository
curl -1sLf https://github.com/rabbitmq/signing-keys/releases/download/3.0/cloudsmith.rabbitmq-server.9F4587F226208342.key | sudo gpg --dearmor | sudo tee /usr/share/keyrings/rabbitmq.9F4587F226208342.gpg > /dev/null

## Add apt repositories maintained by Team RabbitMQ
sudo tee /etc/apt/sources.list.d/rabbitmq.list <<EOF
## Provides modern Erlang/OTP releases
##
deb [arch=amd64 signed-by=/usr/share/keyrings/rabbitmq.E495BB49CC4BBE5B.gpg] https://ppa1.rabbitmq.com/rabbitmq/rabbitmq-erlang/deb/ubuntu noble main
deb-src [signed-by=/usr/share/keyrings/rabbitmq.E495BB49CC4BBE5B.gpg] https://ppa1.rabbitmq.com/rabbitmq/rabbitmq-erlang/deb/ubuntu noble main

# another mirror for redundancy
deb [arch=amd64 signed-by=/usr/share/keyrings/rabbitmq.E495BB49CC4BBE5B.gpg] https://ppa2.rabbitmq.com/rabbitmq/rabbitmq-erlang/deb/ubuntu noble main
deb-src [signed-by=/usr/share/keyrings/rabbitmq.E495BB49CC4BBE5B.gpg] https://ppa2.rabbitmq.com/rabbitmq/rabbitmq-erlang/deb/ubuntu noble main

## Provides RabbitMQ
##
deb [arch=amd64 signed-by=/usr/share/keyrings/rabbitmq.9F4587F226208342.gpg] https://ppa1.rabbitmq.com/rabbitmq/rabbitmq-server/deb/ubuntu noble main
deb-src [signed-by=/usr/share/keyrings/rabbitmq.9F4587F226208342.gpg] https://ppa1.rabbitmq.com/rabbitmq/rabbitmq-server/deb/ubuntu noble main

# another mirror for redundancy
deb [arch=amd64 signed-by=/usr/share/keyrings/rabbitmq.9F4587F226208342.gpg] https://ppa2.rabbitmq.com/rabbitmq/rabbitmq-server/deb/ubuntu noble main
deb-src [signed-by=/usr/share/keyrings/rabbitmq.9F4587F226208342.gpg] https://ppa2.rabbitmq.com/rabbitmq/rabbitmq-server/deb/ubuntu noble main
EOF

## Update package indices
sudo apt-get update -y

## Install Erlang packages
sudo apt-get install -y erlang-base \
                        erlang-asn1 erlang-crypto erlang-eldap erlang-ftp erlang-inets \
                        erlang-mnesia erlang-os-mon erlang-parsetools erlang-public-key \
                        erlang-runtime-tools erlang-snmp erlang-ssl \
                        erlang-syntax-tools erlang-tftp erlang-tools erlang-xmerl

## Install rabbitmq-server and its dependencies
sudo apt-get install rabbitmq-server -y --fix-missing
```

O rabbitMQ já irá manter um server aberto, feche o servidor usando o comando:<br>
`sudo systemctl stop rabbitmq-server`

Com todas as dependências instaladas procure pela pasta ~/go/src. Dentro dessa pasta, extraia a pasta GoChatDist do .zip.<br>
Para executar o programa, dentro da pasta GoChatDist no terminal execute os seguintes comandos (os passos abaixo estão pensados para serem executados no VScode, e precisaram ser adaptados caso não o utilize):<br>
`docker-compose up -d`<br>
Utilizando a porta localhost:15672 você poderá acessar a UI; login e senhar é guest.

Garante que um servidor web local esteja online (no VScode a extensão Live Server faz esse trabalho, uma vez instalada você poderá ligar o servidor web clicando no ícone na lateral inferior direita conforme "Go Live").<br>
Utilize os comandos do makefile:
```bash
make proto
make tidy
make build-cliente
make build-server
make run-server
```

Abra outro terminal na mesma pasta GoChatDist e execute:
```bash
make run-client
```

Após o rodar o client, no navegador use o localhost na porta do Live Server (no terminal do VScode, na aba PORTS será a porta com o processo "Code Extension Host") ou clique com o botão direito do mouse no arquivo web/index.htlm e clique na opção abrir com o Live Server. Para rodar vários cliente, abra outras abas no navegador com o mesmo link do localhost (ex: http://localhost:5500/web/).<br>
Nos terminais será exibido os logs dos servidor e dos clientes enquanto as mensagens são enviadas entre os usuários.