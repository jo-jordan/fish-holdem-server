package domain

type TableDO struct {
	TableID          int64
	PlayerListBySeat []PlayerDO
	TotalPot         int64
	Status           string
	Countdown        int64
	BetRate          string
	CardsOnTable     []string
	PlayerSize       int // config
	CardsNotUsed     []string
}

func (t *TableDO) InitCards() {
	t.CardsNotUsed = []string{
		"101",
		"102",
		"103",
		"104",
		"105",
		"106",
		"107",
		"108",
		"109",
		"110",
		"111",
		"112",
		"113",

		"201",
		"202",
		"203",
		"204",
		"205",
		"206",
		"207",
		"208",
		"209",
		"210",
		"211",
		"212",
		"213",

		"301",
		"302",
		"303",
		"304",
		"305",
		"306",
		"307",
		"308",
		"309",
		"310",
		"311",
		"312",
		"313",

		"401",
		"402",
		"403",
		"404",
		"405",
		"406",
		"407",
		"408",
		"409",
		"410",
		"411",
		"412",
		"413",
	}
}
