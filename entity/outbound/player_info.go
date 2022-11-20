package outbound

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
	PlayerList  []Player `json:"player_list"`
	DataType    string   `json:"data_type"`
	CurPlayerID int64    `json:"cur_player_id"`
}

type Player struct {
	Username    string   `json:"username"`
	ID          int64    `json:"id"`
	Balance     float64  `json:"balance"`
	Bet         int64    `json:"bet"`
	Status      string   `json:"status"`
	Role        string   `json:"role"`
	IsOperator  bool     `json:"is_operator"`
	CardsInHand []string `json:"cards_in_hand"`
}
