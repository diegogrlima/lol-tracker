# Endpoints

## Preparar o ambiente

Execute os comandos a partir da pasta `backend`.

```bash
set -a
source .env
set +a
```

## Champion API

O Champion API utiliza a porta `8082` por padrão e consulta os dados estáticos do
Data Dragon. Esse serviço não precisa da chave da Riot.

### Iniciar o serviço

```bash
go run ./cmd/champion-api
```

### Listar campeões

```bash
curl "http://localhost:8082/champions"
```

Os campeões são retornados em ordem alfabética, incluindo a URL da imagem:

```json
{
  "champions": [
    {
      "id": "Aatrox",
      "key": "266",
      "name": "Aatrox",
      "title": "a Espada Darkin",
      "blurb": "Descrição resumida do campeão",
      "tags": ["Fighter"],
      "imageUrl": "https://ddragon.leagueoflegends.com/cdn/16.1.1/img/champion/Aatrox.png"
    }
  ]
}
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

As consultas usam Redis com TTLs independentes:

```env
MATCH_IDS_CACHE_TTL=5m
MATCH_DETAIL_CACHE_TTL=24h
```

As listas são armazenadas com chaves no formato:

```text
match:ids:{puuid}:{start}:{count}
```

Os detalhes usam:

```text
match:detail:{matchID}
```

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

### Buscar detalhes de uma partida

Copie um dos valores retornados em `matchIds` e substitua `SEU_MATCH_ID`.

```bash
curl "http://localhost:8081/matches/SEU_MATCH_ID"
```

Exemplo:

```bash
curl "http://localhost:8081/matches/BR1_1234567890"
```

A resposta contém os metadados da partida, informações da partida e estatísticas dos participantes:

```json
{
  "match": {
    "metadata": {
      "dataVersion": "2",
      "matchId": "BR1_1234567890",
      "participants": [
        "PUUID_DO_PARTICIPANTE"
      ]
    },
    "info": {
      "gameMode": "CLASSIC",
      "queueId": 420,
      "participants": [
        {
          "puuid": "PUUID_DO_PARTICIPANTE",
          "championName": "Ahri",
          "kills": 10,
          "deaths": 2,
          "assists": 8,
          "win": true
        }
      ]
    }
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

6. Copie um `matchId` da resposta e busque os detalhes:

```bash
curl "http://localhost:8081/matches/SEU_MATCH_ID"
```
