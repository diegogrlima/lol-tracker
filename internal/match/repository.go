package match

import "context"

// Repository contracts will be defined from the needs of the first use case.

type ListOptions struct {
	Start int
	Count int
}

type Match struct {
	Metadata Metadata `json:"metadata"`
	Info     Info     `json:"info"`
}

type Metadata struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type Info struct {
	GameCreation       int64         `json:"gameCreation"`
	GameDuration       int64         `json:"gameDuration"`
	GameEndTimestamp   int64         `json:"gameEndTimestamp"`
	GameMode           string        `json:"gameMode"`
	GameStartTimestamp int64         `json:"gameStartTimestamp"`
	GameType           string        `json:"gameType"`
	GameVersion        string        `json:"gameVersion"`
	MapID              int           `json:"mapId"`
	PlatformID         string        `json:"platformId"`
	QueueID            int           `json:"queueId"`
	Participants       []Participant `json:"participants"`
}

type Participant struct {
	PUUID                       string `json:"puuid"`
	RiotIDGameName              string `json:"riotIdGameName"`
	RiotIDTagLine               string `json:"riotIdTagline"`
	ChampionID                  int    `json:"championId"`
	ChampionName                string `json:"championName"`
	TeamID                      int    `json:"teamId"`
	TeamPosition                string `json:"teamPosition"`
	Kills                       int    `json:"kills"`
	Deaths                      int    `json:"deaths"`
	Assists                     int    `json:"assists"`
	GoldEarned                  int    `json:"goldEarned"`
	TotalDamageDealtToChampions int    `json:"totalDamageDealtToChampions"`
	TotalMinionsKilled          int    `json:"totalMinionsKilled"`
	NeutralMinionsKilled        int    `json:"neutralMinionsKilled"`
	VisionScore                 int    `json:"visionScore"`
	Win                         bool   `json:"win"`
}

type Repository interface {
	ListIDsByPUUID(
		ctx context.Context,
		puuid string,
		options ListOptions,
	) ([]string, error)

	GetByID(ctx context.Context, matchID string) (*Match, error)
}
