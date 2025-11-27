# Auction Labs 3

## 🚀 Como Rodar o Projeto

### Usando Docker Compose (Recomendado)

```bash
docker-compose up --build -d
```

Este comando irá:
- Construir a imagem Docker da aplicação
- Iniciar o container da aplicação na porta `8080`
- Iniciar o MongoDB na porta `27017`

A aplicação estará disponível em: `http://localhost:8080`

### 1. Criar Leilão (POST)
```http
POST http://localhost:8080/auction
Content-Type: application/json

{
  "product_name": "produto teste",
  "category": "teste api rest",
  "description": "teste do teste",
  "condition": 1
}
```

**Resposta:** `201 Created`

### 2. Buscar Leilão por ID (GET)
```http
GET http://localhost:8080/auction/{auctionId}
```

### 3. Listar Leilões (GET)
```http
GET http://localhost:8080/auction?status=0&productName=produto&category=teste
```

**Parâmetros:**
- `status`: Condição do produto (0=Novo, 1=Usado, 2=Recondicionado)
- `productName`: Nome do produto (opcional)
- `category`: Categoria (opcional)

### 4. Criar Lance (POST)
```http
POST http://localhost:8080/bid
Content-Type: application/json

{
  "user_id": "uuid-do-usuario",
  "auction_id": "uuid-do-leilao",
  "amount": 100.50
}
```

### 5. Listar Lances por Leilão (GET)
```http
GET http://localhost:8080/bid/{auctionId}
```

### 6. Buscar Lance Vencedor (GET)
```http
GET http://localhost:8080/auction/winner/{auctionId}
```

### 7. Buscar Usuário (GET)
```http
GET http://localhost:8080/user/{userId}
```

## 📁 Estrutura do Projeto

```
.
├── cmd/auction/               # Ponto de entrada da aplicação
│   ├── main.go               # Arquivo principal
│   └── .env                  # Variáveis de ambiente
├── configuration/            # Configurações
│   ├── database/
│   ├── logger/
│   └── rest_err/
├── internal/
│   ├── entity/               # Entidades do domínio
│   ├── infra/                # Infraestrutura (BD, API)
│   ├── usecase/              # Casos de uso
│   └── internal_error/       # Tratamento de erros
├── api/                      # Exemplos de requisições HTTP
├── docker-compose.yaml       # Orquestração de containers
├── Dockerfile               # Imagem Docker da aplicação
├── go.mod                   # Dependências Go
└── README.md                # Este arquivo
```

## 🛠️ Variáveis de Ambiente

As variáveis estão configuradas em `cmd/auction/.env`:

```env
MONGODB_URL=mongodb://mongoDB:27017
MONGODB_DB=auctions
MAX_BATCH_SIZE=15
MAX_BATCH_SIZE_INTERVAL=6m
MAX_INTERVAL_DURATION_AUCTION=1m
```