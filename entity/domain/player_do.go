package domain

import "github.com/gorilla/websocket"

type PlayerDO struct {
	Conn        *websocket.Conn
	Username    string   `json:"username"`
	ID          int64    `json:"id"`
	Balance     float64  `json:"balance"`
	Bet         int64    `json:"bet"`
	Status      string   `json:"status"`
	Role        string   `json:"role"`
	IsOperator  bool     `json:"is_operator"`
	CardsInHand []string `json:"cards_in_hand"`
}
