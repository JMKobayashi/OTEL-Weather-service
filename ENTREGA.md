# 🎯 Entrega - OTEL Weather Service

## 📋 Resumo da Implementação

Foi implementado com sucesso um sistema distribuído em Go que atende todos os requisitos solicitados:

### ✅ Objetivo Principal
Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade, implementando OTEL (OpenTelemetry) e Zipkin.

## 🏗️ Arquitetura Implementada

### Serviço A (CEP Validator) - Porta 8080
**Responsabilidade:** Validação de CEP e orquestração

**Endpoint:** `POST /cep`
```json
{
    "cep": "29902555"
}
```

**Funcionalidades:**
- ✅ Validação de CEP (8 dígitos)
- ✅ Encaminhamento para Serviço B via HTTP
- ✅ Tratamento de erros adequado
- ✅ Tracing distribuído com OTEL

### Serviço B (Weather Service) - Porta 8081
**Responsabilidade:** Consulta de localização e temperatura

**Endpoint:** `GET /weather/:zipcode`

**Funcionalidades:**
- ✅ Consulta ViaCEP para localização
- ✅ Consulta WeatherAPI para temperatura
- ✅ Conversão de temperaturas (Celsius, Fahrenheit, Kelvin)
- ✅ Tracing detalhado de cada operação

## 📡 APIs Externas Utilizadas

### ViaCEP API
- **URL:** https://viacep.com.br/
- **Uso:** Consulta de localização por CEP
- **Resposta:** Dados da localidade (cidade, estado, etc.)

### WeatherAPI
- **URL:** https://www.weatherapi.com/
- **Uso:** Consulta de temperatura atual
- **Resposta:** Temperatura em Celsius

## 🌡️ Conversões de Temperatura

Implementadas conforme especificação:
- **Fahrenheit:** `F = C * 1.8 + 32`
- **Kelvin:** `K = C + 273`

## 📋 Requisitos Atendidos

### ✅ Serviço A
- [x] Recebe input de 8 dígitos via POST
- [x] Valida se o input é válido (8 dígitos, string)
- [x] Encaminha para Serviço B via HTTP
- [x] Retorna erro 422 para CEP inválido
- [x] Implementa OTEL com spans

### ✅ Serviço B
- [x] Recebe CEP válido de 8 dígitos
- [x] Pesquisa CEP e encontra localização
- [x] Retorna temperaturas em Celsius, Fahrenheit e Kelvin
- [x] Retorna código 200 em caso de sucesso
- [x] Retorna código 422 para CEP inválido
- [x] Retorna código 404 para CEP não encontrado
- [x] Implementa OTEL com spans detalhados

### ✅ OTEL + Zipkin
- [x] Tracing distribuído entre Serviço A e Serviço B
- [x] Spans para medir tempo de resposta
- [x] Spans para busca de CEP
- [x] Spans para busca de temperatura
- [x] Visualização no Zipkin

## 🔍 Observabilidade Implementada

### OpenTelemetry (OTEL)
- **Configuração:** OTEL Collector com gRPC
- **Spans implementados:**
  - `validate_cep` (Serviço A)
  - `call_weather_service` (Serviço A)
  - `validate_zipcode` (Serviço B)
  - `get_location_by_zipcode` (Serviço B)
  - `get_temperature_by_location` (Serviço B)
  - `viacep_request` (Serviço B)
  - `weatherapi_request` (Serviço B)

### Zipkin
- **URL:** http://localhost:9411
- **Funcionalidade:** Visualização de traces distribuídos
- **Integração:** OTEL Collector envia traces para Zipkin

## 🐳 Containerização

### Docker Compose
- **OTEL Collector:** Configurado para receber traces e enviar para Zipkin
- **Zipkin:** Interface web para visualização de traces
- **Serviço A:** Containerizado com configuração OTEL
- **Serviço B:** Containerizado com configuração OTEL

### Configurações
- **Rede:** Bridge network para comunicação entre serviços
- **Portas:** 8080 (A), 8081 (B), 9411 (Zipkin), 4317 (OTEL)
- **Variáveis de ambiente:** Configuráveis via `.env`

## 🚀 Como Executar

### 1. Configuração Inicial
```bash
# Clone o repositório
git clone <repository-url>
cd otel-weather-service

# Configure as variáveis de ambiente
cp env.example .env
# Edite .env e adicione sua WEATHER_API_KEY
```

### 2. Executar com Docker Compose
```bash
# Build e execução completa
docker-compose up --build
```

### 3. Serviços Disponíveis
- **Serviço A:** http://localhost:8080
- **Serviço B:** http://localhost:8081
- **Zipkin:** http://localhost:9411

## 🧪 Testes

### Teste do Serviço A:
```bash
# Health check
curl http://localhost:8080/health

# Teste com CEP válido
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310900"}'

# Teste com CEP inválido
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "123"}'
```

### Teste do Serviço B diretamente:
```bash
# Health check
curl http://localhost:8081/health

# Teste com CEP válido
curl http://localhost:8081/weather/01310900

# Teste com CEP inválido
curl http://localhost:8081/weather/123

# Teste com CEP inexistente
curl http://localhost:8081/weather/99999999
```

## 📊 Respostas Esperadas

### ✅ Sucesso (200):
```json
{
    "city": "São Paulo",
    "temp_C": 28.5,
    "temp_F": 83.3,
    "temp_K": 301.65
}
```

### ❌ CEP inválido (422):
```json
{
    "error": "invalid zipcode"
}
```

### ❌ CEP não encontrado (404):
```json
{
    "error": "can not find zipcode"
}
```

## 📁 Estrutura do Projeto

```
otel-weather-service/
├── service-a/                 # Serviço A (CEP Validator)
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handlers/
│   │   │   └── cep_handler.go
│   │   ├── services/
│   │   │   └── cep_service.go
│   │   └── models/
│   │       └── cep.go
│   ├── Dockerfile
│   └── go.mod
├── service-b/                 # Serviço B (Weather Service)
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handlers/
│   │   │   └── weather_handler.go
│   │   ├── services/
│   │   │   └── weather_service.go
│   │   └── models/
│   │       └── weather.go
│   ├── Dockerfile
│   └── go.mod
├── docker-compose.yml
├── otel-collector-config.yaml
├── env.example
├── test.http
├── README.md
├── DEPLOY_INSTRUCTIONS.md
├── IMPLEMENTATION_SUMMARY.md
└── ENTREGA.md
```

## 📝 Tecnologias Utilizadas

- **Go 1.21:** Linguagem principal
- **Gin:** Framework web
- **OpenTelemetry:** Instrumentação de telemetria
- **Zipkin:** Visualização de traces
- **Docker:** Containerização
- **Docker Compose:** Orquestração
- **ViaCEP API:** Consulta de CEP
- **WeatherAPI:** Consulta de temperatura

## 🎉 Resultado Final

O sistema implementa **completamente** todos os requisitos solicitados:

- ✅ **Dois serviços distribuídos** (A e B)
- ✅ **Validação adequada de CEP** (8 dígitos)
- ✅ **Consulta de temperatura** com conversões (Celsius, Fahrenheit, Kelvin)
- ✅ **Tratamento de erros** conforme especificação (422, 404)
- ✅ **Tracing distribuído** com OTEL
- ✅ **Visualização no Zipkin**
- ✅ **Containerização completa** com Docker Compose
- ✅ **Documentação detalhada** de uso e deploy

## 🔗 Links Úteis

- **Zipkin:** http://localhost:9411 (após execução)
- **WeatherAPI:** https://www.weatherapi.com/ (para obter chave da API)
- **ViaCEP:** https://viacep.com.br/ (API de consulta de CEP)

---

**Status:** ✅ **IMPLEMENTAÇÃO COMPLETA**

Todos os requisitos foram atendidos e o sistema está pronto para uso em ambiente de desenvolvimento. 