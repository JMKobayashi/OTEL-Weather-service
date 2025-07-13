# Resumo da Implementação - OTEL Weather Service

## 🎯 Objetivo Alcançado

Foi implementado com sucesso um sistema distribuído em Go que recebe um CEP, identifica a cidade e retorna o clima atual com temperaturas em Celsius, Fahrenheit e Kelvin, implementando OTEL (OpenTelemetry) e Zipkin para tracing distribuído.

## 🏗️ Arquitetura Implementada

### Serviço A (CEP Validator) - Porta 8080
- **Responsabilidade:** Validação de CEP e orquestração
- **Endpoint:** `POST /cep`
- **Funcionalidades:**
  - Validação de CEP (8 dígitos)
  - Encaminhamento para Serviço B
  - Tratamento de erros adequado
  - Tracing distribuído

### Serviço B (Weather Service) - Porta 8081
- **Responsabilidade:** Consulta de localização e temperatura
- **Endpoint:** `GET /weather/:zipcode`
- **Funcionalidades:**
  - Consulta ViaCEP para localização
  - Consulta WeatherAPI para temperatura
  - Conversão de temperaturas (Celsius, Fahrenheit, Kelvin)
  - Tracing detalhado de cada operação

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

## 🧪 Testes

### Endpoints Testados
- **Health checks:** Ambos os serviços
- **CEP válido:** Retorna temperaturas corretas
- **CEP inválido:** Retorna erro 422
- **CEP inexistente:** Retorna erro 404

### Observabilidade Testada
- **Traces distribuídos:** Visualizáveis no Zipkin
- **Spans detalhados:** Cada operação rastreada
- **Tempo de resposta:** Medido para cada operação

## 📊 Métricas e Monitoramento

### Disponíveis no Zipkin
- Tempo de resposta de cada operação
- Traces distribuídos entre serviços
- Detalhes das chamadas HTTP externas
- Erros e exceções

### Spans Implementados
1. **Serviço A:**
   - Validação de CEP
   - Chamada HTTP para Serviço B

2. **Serviço B:**
   - Validação de CEP
   - Consulta ViaCEP
   - Consulta WeatherAPI
   - Conversão de temperaturas

## 🚀 Deploy

### Ambiente de Desenvolvimento
```bash
# Configurar variáveis
cp env.example .env
# Editar .env com WEATHER_API_KEY

# Executar
docker-compose up --build
```

### Serviços Disponíveis
- **Serviço A:** http://localhost:8080
- **Serviço B:** http://localhost:8081
- **Zipkin:** http://localhost:9411

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

O sistema implementa completamente os requisitos solicitados:
- ✅ Dois serviços distribuídos
- ✅ Validação adequada de CEP
- ✅ Consulta de temperatura com conversões
- ✅ Tratamento de erros conforme especificação
- ✅ Tracing distribuído com OTEL
- ✅ Visualização no Zipkin
- ✅ Containerização completa
- ✅ Documentação detalhada 