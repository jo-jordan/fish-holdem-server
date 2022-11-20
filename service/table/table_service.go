package table_service

import (
	"github.com/jo-jordan/fish-holdem-server/entity/domain"
	"github.com/jo-jordan/fish-holdem-server/entity/outbound"
	"github.com/jo-jordan/fish-holdem-server/misc"
	"github.com/jo-jordan/fish-holdem-server/misc/global"
	"github.com/jo-jordan/fish-holdem-server/util"
	"math/rand"
	"time"
)

func MatchTable(playerDO *domain.PlayerDO) (*outbound.TableInfo, []*outbound.PlayerInfo) {

	var tableJoined domain.TableDO

	tableNum := 0
	joined := false
	global.TableMap.Range(func(key, value any) bool {
		if value == nil {
			return false
		}
		tableDO := value.(domain.TableDO)
		tableNum += 1

		if len(tableDO.PlayerListBySeat) < tableDO.PlayerSize {
			tableDO.PlayerListBySeat = append(tableDO.PlayerListBySeat, *playerDO)
			joined = true
			tableJoined = tableDO
			global.TableMap.Store(tableJoined.TableID, tableJoined)
			return false
		}

		return true
	})

	if !joined {
		tableJoined = createTable()
		tableJoined.PlayerListBySeat = append(tableJoined.PlayerListBySeat, *playerDO)
		global.TableMap.Store(tableJoined.TableID, tableJoined)
	}

	tableInfo := outbound.TableInfo{
		TableID:      tableJoined.TableID,
		TotalPot:     tableJoined.TotalPot,
		Status:       tableJoined.Status,
		Countdown:    tableJoined.Countdown,
		BetRate:      tableJoined.BetRate,
		CardsOnTable: tableJoined.CardsOnTable,
		PlayerSize:   tableJoined.PlayerSize,
		DataType:     misc.DTTableInfo,
	}

	if len(tableJoined.PlayerListBySeat) == tableJoined.PlayerSize {
		deal(&tableJoined)
	}

	playerInfos := make([]*outbound.PlayerInfo, 0)

	for _, v := range tableJoined.PlayerListBySeat {
		playerList := make([]outbound.Player, 0)

		for _, p := range tableJoined.PlayerListBySeat {
			cih := p.CardsInHand
			if v.ID != p.ID {
				cih = []string{}
			}

			playerList = append(playerList,
				outbound.Player{
					ID:          p.ID,
					Username:    p.Username,
					Balance:     p.Balance,
					Bet:         p.Bet,
					Status:      p.Status,
					Role:        p.Role,
					IsOperator:  p.IsOperator,
					CardsInHand: cih,
				},
			)
		}

		pi := outbound.PlayerInfo{
			CurPlayerID: v.ID,
			DataType:    misc.DTPlayerInfo,
			PlayerList:  playerList,
		}

		playerInfos = append(playerInfos, &pi)
	}

	return &tableInfo, playerInfos
}

func deal(table *domain.TableDO) {
	table.InitCards()

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(table.CardsNotUsed), func(i, j int) {
		table.CardsNotUsed[i], table.CardsNotUsed[j] = table.CardsNotUsed[j], table.CardsNotUsed[i]
	})

	for i, p := range table.PlayerListBySeat {
		p.CardsInHand = []string{}
		p.CardsInHand = append(p.CardsInHand, table.CardsNotUsed[0])
		table.CardsNotUsed = table.CardsNotUsed[1:]

		p.CardsInHand = append(p.CardsInHand, table.CardsNotUsed[0])
		table.CardsNotUsed = table.CardsNotUsed[1:]

		table.PlayerListBySeat[i] = p
	}
	global.TableMap.Store(table.TableID, *table)
}

func createTable() domain.TableDO {
	return domain.TableDO{
		TableID:          util.GenID(),
		TotalPot:         0,
		Status:           "",
		Countdown:        20,
		BetRate:          "10/20",
		PlayerSize:       2,
		CardsOnTable:     []string{},
		PlayerListBySeat: make([]domain.PlayerDO, 0),
	}
}
