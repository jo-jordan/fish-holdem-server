package domain

type TableDO struct {
	TableID          int64 `json:"table_id"`
	PlayerListBySeat []PlayerDO
	TotalPot         int64    `json:"total_pot"`
	Status           string   `json:"status"`
	Countdown        int64    `json:"countdown"`
	BetRate          string   `json:"bet_rate"`
	CardsOnTable     []string `json:"cards_on_table"`
	PlayerSize       int      `json:"player_size"` // config
}
