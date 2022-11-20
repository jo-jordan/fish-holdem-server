package domain

import (
	"sync/atomic"
	"time"
)

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
	CurActionPlayer  *PlayerDO
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

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

func (t *TableDO) StartLoop() {
	ticker := time.NewTicker(time.Duration(t.Countdown) * time.Second)

	defer ticker.Stop()

	var next uint32 = 0
	for {
		select {
		case <-ticker.C:
			nextPlayer := t.PlayerListBySeat[next]
			t.CurActionPlayer = &nextPlayer
			atomic.AddUint32(&next, 1)
			if atomic.LoadUint32(&next) == uint32(len(t.PlayerListBySeat)-1) {
				atomic.StoreUint32(&next, 0)
			}
		}
	}
}
