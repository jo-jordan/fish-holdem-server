package outbound

import "encoding/json"

func UnmarshalGameInfo(data []byte) (TableInfo, error) {
	var r TableInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TableInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type TableInfo struct {
	TableID      int64    `json:"table_id"`
	TotalPot     int64    `json:"total_pot"`
	Status       string   `json:"status"`
	Countdown    int64    `json:"countdown"`
	BetRate      string   `json:"bet_rate"`
	CardsOnTable []string `json:"cards_on_table"`
	PlayerSize   int64    `json:"player_size"`
	DataType     string   `json:"data_type"`
}
