# Instruções de Deploy - OTEL Weather Service

## 📋 Pré-requisitos

1. **Docker e Docker Compose** instalados
2. **Chave da API WeatherAPI** (gratuita em https://www.weatherapi.com/)

## 🚀 Executando o Projeto

### 1. Configuração Inicial

1. Clone o repositório (se ainda não fez):
```bash
git clone <repository-url>
cd otel-weather-service
```

2. Configure as variáveis de ambiente:
```bash
cp env.example .env
```

3. Edite o arquivo `.env` e adicione sua chave da API:
```bash
WEATHER_API_KEY=sua_chave_aqui
```

### 2. Executando com Docker Compose

```bash
# Build e execução completa
docker-compose up --build

# Ou em background
docker-compose up -d --build
```

### 3. Verificando os Serviços

Após a execução, os seguintes serviços estarão disponíveis:

- **Serviço A (CEP Validator):** http://localhost:8080
- **Serviço B (Weather Service):** http://localhost:8081
- **Zipkin (Tracing):** http://localhost:9411
- **OTEL Collector:** http://localhost:13133 (health check)

### 4. Testando os Endpoints

#### Teste do Serviço A:
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

#### Teste do Serviço B diretamente:
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

## 🔍 Observabilidade

### Zipkin
Acesse http://localhost:9411 para visualizar os traces distribuídos.

**Funcionalidades:**
- Visualização de traces entre Serviço A e Serviço B
- Tempo de resposta de cada operação
- Detalhes das chamadas HTTP para APIs externas (ViaCEP e WeatherAPI)

### OTEL Collector
O collector está configurado para:
- Receber traces via gRPC na porta 4317
- Processar e enviar para Zipkin
- Health check disponível em http://localhost:13133

## 🐳 Comandos Docker Úteis

```bash
# Ver logs de todos os serviços
docker-compose logs

# Ver logs de um serviço específico
docker-compose logs service-a
docker-compose logs service-b

# Parar todos os serviços
docker-compose down

# Rebuild e restart
docker-compose down
docker-compose up --build

# Executar apenas os serviços (sem observabilidade)
docker-compose up service-a service-b
```

## 🧪 Testes Automatizados

Use o arquivo `test.http` para testar todos os endpoints:

1. Abra o arquivo no VS Code ou sua IDE
2. Instale a extensão "REST Client" se necessário
3. Execute os testes clicando em "Send Request"

## 🔧 Troubleshooting

### Problemas Comuns

1. **Erro de conexão com OTEL Collector:**
   - Verifique se o collector está rodando: `docker-compose logs otel-collector`
   - Aguarde alguns segundos para o collector inicializar

2. **Erro de WeatherAPI:**
   - Verifique se a chave da API está correta no arquivo `.env`
   - Teste a chave diretamente: `curl "http://api.weatherapi.com/v1/current.json?key=SUA_CHAVE&q=Sao Paulo"`

3. **Serviços não iniciam:**
   - Verifique os logs: `docker-compose logs`
   - Verifique se as portas não estão em uso: `netstat -tulpn | grep :808`

### Logs Úteis

```bash
# Ver logs em tempo real
docker-compose logs -f

# Ver logs de um serviço específico
docker-compose logs -f service-a

# Ver logs do OTEL Collector
docker-compose logs -f otel-collector
```

## 📊 Monitoramento

### Métricas Disponíveis

- **Tempo de resposta** de cada operação
- **Taxa de erro** por endpoint
- **Tempo de resposta** das APIs externas (ViaCEP e WeatherAPI)

### Visualização no Zipkin

1. Acesse http://localhost:9411
2. Clique em "Find Traces"
3. Selecione o serviço desejado
4. Visualize os traces e spans

## 🚀 Deploy em Produção

Para deploy em produção, considere:

1. **Variáveis de ambiente** para configurações sensíveis
2. **Health checks** para monitoramento
3. **Logs estruturados** para melhor observabilidade
4. **Métricas** para monitoramento de performance
5. **Backup** das configurações do OTEL Collector

## 📝 Notas Importantes

- O projeto usa **tracing distribuído** com OpenTelemetry
- Todos os **spans** são enviados para o Zipkin
- As **conversões de temperatura** são feitas localmente
- O **contexto de tracing** é propagado entre os serviços
- **Graceful shutdown** está implementado em ambos os serviços 