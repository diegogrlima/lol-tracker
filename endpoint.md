# Endpoints

## Preparar o ambiente

Execute os comandos a partir da pasta `backend`.

```bash
set -a
source .env
set +a
```

## Player API

O Player API utiliza a porta `8080` por padrão.

### Iniciar o serviço

```bash
go run ./cmd/player-api
```

### Health check

```bash
curl "http://localhost:8080/health"
```

Resposta esperada:

```json
{
  "health": {
    "message": "API funcionando",
    "status": true
  }
}
```

### Buscar jogador pelo Riot ID

Substitua `NOME_DO_JOGADOR` e `TAG`.

```bash
curl "http://localhost:8080/players/NOME_DO_JOGADOR/TAG"
```

Exemplo:

```bash
curl "http://localhost:8080/players/Player/BR1"
```

Resposta esperada:

```json
{
  "player": {
    "puuid": "PUUID_RETORNADO_PELA_RIOT",
    "gameName": "Player",
    "tagLine": "BR1"
  }
}
```

Copie o valor de `puuid` da resposta para testar os endpoints do Match API.

## Match API

O Match API utiliza a porta `8081` por padrão.

### Iniciar o serviço

```bash
go run ./cmd/match-api
```

### Health check

```bash
curl "http://localhost:8081/health"
```

### Listar as últimas 20 partidas

Substitua somente `SEU_PUUID`.

```bash
curl "http://localhost:8081/matches/by-puuid/SEU_PUUID"
```

Esse comando utiliza os valores padrão:

```text
start=0
count=20
```

### Listar partidas com paginação explícita

```bash
curl "http://localhost:8081/matches/by-puuid/SEU_PUUID?start=0&count=20"
```

### Listar as primeiras cinco partidas

```bash
curl "http://localhost:8081/matches/by-puuid/SEU_PUUID?start=0&count=5"
```

### Listar cinco partidas a partir da décima

```bash
curl "http://localhost:8081/matches/by-puuid/SEU_PUUID?start=10&count=5"
```

### Listar o máximo de 100 partidas

```bash
curl "http://localhost:8081/matches/by-puuid/SEU_PUUID?start=0&count=100"
```

Resposta esperada:

```json
{
  "matchIds": [
    "BR1_1234567890",
    "BR1_9876543210"
  ],
  "pagination": {
    "start": 0,
    "count": 20
  }
}
```

## Testes de validação

### `start` negativo

```bash
curl -i "http://localhost:8081/matches/by-puuid/SEU_PUUID?start=-1&count=20"
```

Resposta esperada: `400 Bad Request`.

### `count` igual a zero

```bash
curl -i "http://localhost:8081/matches/by-puuid/SEU_PUUID?start=0&count=0"
```

Resposta esperada: `400 Bad Request`.

### `count` acima do limite

```bash
curl -i "http://localhost:8081/matches/by-puuid/SEU_PUUID?start=0&count=101"
```

Resposta esperada: `400 Bad Request`.

### Paginação não numérica

```bash
curl -i "http://localhost:8081/matches/by-puuid/SEU_PUUID?start=abc&count=20"
```

Resposta esperada: `400 Bad Request`.

## Fluxo completo

1. Inicie o Redis:

```bash
docker compose up -d redis
```

2. Inicie o Player API em um terminal:

```bash
set -a; source .env; set +a; go run ./cmd/player-api
```

3. Inicie o Match API em outro terminal:

```bash
set -a; source .env; set +a; go run ./cmd/match-api
```

4. Busque um jogador e copie o `puuid`:

```bash
curl "http://localhost:8080/players/NOME_DO_JOGADOR/TAG"
```

5. Substitua `SEU_PUUID` e liste as partidas:

```bash
curl "http://localhost:8081/matches/by-puuid/SEU_PUUID?start=0&count=20"
```
