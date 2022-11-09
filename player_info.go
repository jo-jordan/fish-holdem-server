package main

import "encoding/json"

func UnmarshalPlayerInfo(data []byte) (PlayerInfo, error) {
	var r PlayerInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *PlayerInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type PlayerInfo struct {
	PlayerList []PlayerList `json:"player_list"`
	DataType   string       `json:"data_type"`
}

type PlayerList struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Balance     float64  `json:"balance"`
	Bet         int64    `json:"bet"`
	Status      string   `json:"status"`
	Role        string   `json:"role"`
	IsOperator  bool     `json:"is_operator"`
	CardsInHand []string `json:"cards_in_hand"`
}
