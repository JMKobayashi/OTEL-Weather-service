# OTEL Weather Service

Sistema distribuído em Go que implementa dois serviços para consulta de temperatura por CEP, com tracing distribuído usando OpenTelemetry e Zipkin.

## 🏗️ Arquitetura

O sistema é composto por dois serviços:

### Serviço A (CEP Validator)
- **Porta:** 8080
- **Responsabilidade:** Validação de CEP e encaminhamento para o Serviço B
- **Endpoint:** `POST /cep`

### Serviço B (Weather Service)
- **Porta:** 8081
- **Responsabilidade:** Consulta de localização e temperatura
- **Endpoint:** `GET /weather/:zipcode`

## 🚀 Funcionalidades

- ✅ Validação de CEP (8 dígitos)
- ✅ Consulta de localização via ViaCEP API
- ✅ Consulta de temperatura via WeatherAPI
- ✅ Conversão automática de temperaturas (Celsius, Fahrenheit, Kelvin)
- ✅ Tracing distribuído com OpenTelemetry
- ✅ Visualização de traces no Zipkin
- ✅ Tratamento adequado de erros
- ✅ Docker Compose para ambiente completo

## 📋 Requisitos

- Go 1.21 ou superior
- Docker e Docker Compose
- Conta no WeatherAPI (https://www.weatherapi.com/)

## 🛠️ Configuração

1. Clone o repositório
2. Copie o arquivo de exemplo de variáveis de ambiente:
```bash
cp .env.example .env
```

3. Edite o arquivo `.env` e adicione sua chave da API do WeatherAPI:
```bash
WEATHER_API_KEY=sua_chave_aqui
```

## 🏃‍♂️ Executando o Projeto

### Com Docker Compose (Recomendado):
```bash
docker-compose up --build
```

### Executando localmente:

#### Serviço A:
```bash
cd service-a
go run cmd/main.go
```

#### Serviço B:
```bash
cd service-b
go run cmd/main.go
```

## 📡 Endpoints

### Serviço A - POST /cep
Recebe um CEP e valida antes de encaminhar para o Serviço B.

**Request:**
```json
{
    "cep": "29902555"
}
```

**Respostas:**

✅ Sucesso (200):
```json
{
    "city": "São Paulo",
    "temp_C": 28.5,
    "temp_F": 83.3,
    "temp_K": 301.65
}
```

❌ CEP inválido (422):
```json
{
    "error": "invalid zipcode"
}
```

### Serviço B - GET /weather/:zipcode
Retorna a temperatura atual para um CEP específico.

**Exemplos:**
```bash
curl http://localhost:8081/weather/01310900
```

**Respostas:**

✅ Sucesso (200):
```json
{
    "city": "São Paulo",
    "temp_C": 28.5,
    "temp_F": 83.3,
    "temp_K": 301.65
}
```

❌ CEP inválido (422):
```json
{
    "error": "invalid zipcode"
}
```

❌ CEP não encontrado (404):
```json
{
    "error": "can not find zipcode"
}
```

## 🔍 Observabilidade

### Zipkin
- **URL:** http://localhost:9411
- **Funcionalidade:** Visualização de traces distribuídos

### OTEL Collector
- **Porta:** 4317 (gRPC)
- **Funcionalidade:** Coleta e processamento de telemetria

## 🧪 Testes

### Testando o Serviço A:
```bash
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310900"}'
```

### Testando o Serviço B diretamente:
```bash
curl http://localhost:8081/weather/01310900
```

## 📊 Tracing

O sistema implementa tracing distribuído com os seguintes spans:

1. **Serviço A:**
   - Validação de CEP
   - Chamada HTTP para Serviço B

2. **Serviço B:**
   - Consulta ViaCEP
   - Consulta WeatherAPI
   - Conversão de temperaturas

## 🐳 Docker

### Build das imagens:
```bash
docker-compose build
```

### Executar apenas os serviços:
```bash
docker-compose up service-a service-b
```

### Executar com observabilidade:
```bash
docker-compose up
```

## 📁 Estrutura do Projeto

```
otel-weather-service/
├── service-a/                 # Serviço A (CEP Validator)
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handlers/
│   │   ├── services/
│   │   └── models/
│   └── Dockerfile
├── service-b/                 # Serviço B (Weather Service)
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handlers/
│   │   ├── services/
│   │   └── models/
│   └── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```
