package main

import "encoding/json"

func UnmarshalGameInfo(data []byte) (GameInfo, error) {
	var r GameInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GameInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GameInfo struct {
	ID           int64    `json:"id"`
	TableID      int64    `json:"table_id"`
	TotalPot     int64    `json:"total_pot"`
	Status       string   `json:"status"`
	Countdown    int64    `json:"countdown"`
	BetRate      string   `json:"bet_rate"`
	CardsOnTable []string `json:"cards_on_table"`
	DataType     string   `json:"data_type"`
}
