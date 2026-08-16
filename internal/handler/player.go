package handler

import (
	"fmt"
	"net/http"
)

type Player struct{}

func (p *Player) GetPlayer(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Busca um player pelo seu Riot ID")
}
