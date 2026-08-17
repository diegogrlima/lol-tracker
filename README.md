# LoL Tracker API

Backend em Go para consultar jogadores, histórico de partidas, campeões e itens de **League of Legends**. O projeto integra a API da Riot Games e o Data Dragon, usando Redis para reduzir chamadas repetidas e melhorar o tempo de resposta.

## Funcionalidades

- busca de jogadores pelo Riot ID (`gameName` e `tagLine`);
- listagem paginada das partidas de um jogador;
- consulta dos detalhes e estatísticas de uma partida;
- catálogo de campeões e itens fornecido pelo Data Dragon;
- cache Redis com TTL independente para jogadores, listas e detalhes de partidas;
- health check, logs estruturados, timeout de requisição e encerramento gracioso;
- testes unitários para serviços, handlers, configuração e adaptadores.

## Arquitetura

A aplicação é formada por três serviços HTTP independentes:

| Serviço | Porta padrão | Responsabilidade |
| --- | ---: | --- |
| Player API | `8080` | Consulta de jogadores na Riot API |
| Match API | `8081` | Histórico e detalhes de partidas |
| Game Data API | `8082` | Campeões e itens do Data Dragon |

```text
Cliente
  ├── Player API ── Riot Account API
  │      └── Redis
  ├── Match API ─── Riot Match API
  │      └── Redis
  └── Game Data API ── Data Dragon
```

O código segue uma separação por domínio, com handlers HTTP, serviços, contratos de repositório e adaptadores externos.

## Tecnologias

- Go 1.26.5
- Chi Router
- Redis 8
- Docker Compose
- Riot Games API
- Data Dragon

## Pré-requisitos

- Go 1.26.5 ou compatível;
- Docker com Docker Compose;
- uma chave de desenvolvimento da Riot Games para as APIs de jogadores e partidas.

O Game Data API pode ser executado sem Redis e sem uma chave da Riot.

## Configuração

Clone o repositório e entre na pasta do backend:

```bash
git clone https://github.com/diegogrlima/lol-tracker.git
cd lol-tracker
```

Crie o arquivo de ambiente a partir do exemplo:

```bash
cp .env.example .env
```

Preencha pelo menos `RIOT_API_KEY` para executar o Player API e o Match API:

```env
RIOT_API_KEY=RGAPI-sua-chave
```

Variáveis disponíveis:

| Variável | Padrão | Descrição |
| --- | --- | --- |
| `SERVER_PORT` | `8080` | Porta do Player API |
| `MATCH_SERVER_PORT` | `8081` | Porta do Match API |
| `GAME_DATA_SERVER_PORT` | `8082` | Porta do Game Data API |
| `REDIS_ADDRESS` | `localhost:6379` | Endereço do Redis |
| `REDIS_PASSWORD` | vazio | Senha do Redis |
| `REDIS_DB` | `0` | Banco lógico do Redis |
| `RIOT_API_KEY` | — | Chave obrigatória para os serviços Riot |
| `RIOT_REGION` | `americas` | Rota regional da Riot |
| `CACHE_TTL` | `5m` | Cache de jogadores |
| `MATCH_IDS_CACHE_TTL` | `5m` | Cache das listas de partidas |
| `MATCH_DETAIL_CACHE_TTL` | `24h` | Cache dos detalhes de partidas |
| `DATA_DRAGON_BASE_URL` | URL oficial | Endereço base do Data Dragon |
| `DATA_DRAGON_VERSION` | `16.1.1` | Versão dos dados estáticos |
| `DATA_DRAGON_LOCALE` | `pt_BR` | Localidade dos dados estáticos |

## Como executar

Inicie o Redis:

```bash
docker compose up -d redis
```

Exporte as variáveis do `.env` no terminal:

```bash
set -a
source .env
set +a
```

Depois, execute cada serviço em um terminal separado:

```bash
go run ./cmd/player-api
```

```bash
go run ./cmd/match-api
```

```bash
go run ./cmd/game-data-api
```

Confirme que os serviços estão ativos:

```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
```

## Endpoints

### Player API

| Método | Endpoint | Descrição |
| --- | --- | --- |
| `GET` | `/health` | Verifica a saúde do serviço |
| `GET` | `/players/{gameName}/{tagLine}` | Busca um jogador pelo Riot ID |

```bash
curl "http://localhost:8080/players/NOME_DO_JOGADOR/BR1"
```

### Match API

| Método | Endpoint | Descrição |
| --- | --- | --- |
| `GET` | `/health` | Verifica a saúde do serviço |
| `GET` | `/matches/by-puuid/{puuid}` | Lista IDs de partidas |
| `GET` | `/matches/{matchID}` | Retorna os detalhes de uma partida |

A listagem aceita `start` (padrão `0`) e `count` (padrão `20`, máximo `100`):

```bash
curl "http://localhost:8081/matches/by-puuid/SEU_PUUID?start=0&count=20"
curl "http://localhost:8081/matches/BR1_1234567890"
```

### Game Data API

| Método | Endpoint | Descrição |
| --- | --- | --- |
| `GET` | `/health` | Verifica a saúde do serviço |
| `GET` | `/champions` | Lista os campeões em ordem alfabética |
| `GET` | `/items` | Lista os itens compráveis da loja |
| `GET` | `/items/{itemID}` | Busca um item pelo ID |

```bash
curl http://localhost:8082/champions
curl http://localhost:8082/items
curl http://localhost:8082/items/1001
```

Consulte [endpoints.md](endpoints.md) para exemplos completos de respostas, paginação e validações.

## Fluxo de uso

1. Consulte `/players/{gameName}/{tagLine}` e copie o `puuid` retornado.
2. Use o PUUID em `/matches/by-puuid/{puuid}` para obter os IDs das partidas.
3. Consulte `/matches/{matchID}` para visualizar os detalhes de uma partida.

## Testes

Execute toda a suíte a partir da raiz do projeto:

```bash
go test ./...
```

Para exibir a cobertura:

```bash
go test -cover ./...
```

## Estrutura do projeto

```text
.
├── cmd/
│   ├── game-data-api/    # inicializa campeões e itens
│   ├── match-api/        # inicializa partidas
│   └── player-api/       # inicializa jogadores
├── internal/
│   ├── champion/        # domínio de campeões
│   ├── config/          # configuração por ambiente
│   ├── item/            # domínio de itens
│   ├── match/           # domínio de partidas
│   ├── platform/        # clientes Riot, Data Dragon e Redis
│   ├── player/          # domínio de jogadores
│   └── server/          # servidor e roteamento HTTP
├── .env.example
├── compose.yml
└── endpoints.md
```

## Licença

Este projeto ainda não possui uma licença definida.
