package table_service

import (
	"github.com/jo-jordan/fish-holdem-server/entity/domain"
	"github.com/jo-jordan/fish-holdem-server/entity/outbound"
	"github.com/jo-jordan/fish-holdem-server/misc/global"
	"github.com/jo-jordan/fish-holdem-server/util"
)

func MatchTable(playerDO *domain.PlayerDO) (*outbound.TableInfo, *outbound.PlayerInfo) {

	var tableJoined domain.TableDO

	tableNum := 0
	joined := false
	global.TableMap.Range(func(key, value any) bool {
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
		DataType:     "table_info",
	}

	playerList := make([]outbound.Player, 0)

	for _, v := range tableJoined.PlayerListBySeat {
		playerList = append(playerList,
			outbound.Player{
				ID:          v.ID,
				Username:    v.Username,
				Balance:     v.Balance,
				Bet:         v.Bet,
				Status:      v.Status,
				Role:        v.Role,
				IsOperator:  v.IsOperator,
				CardsInHand: v.CardsInHand,
			},
		)
	}

	playerInfo := outbound.PlayerInfo{
		DataType:   "player_info",
		PlayerList: playerList,
	}

	return &tableInfo, &playerInfo
}

func createTable() domain.TableDO {
	return domain.TableDO{
		TableID:          util.GenID(),
		TotalPot:         0,
		Status:           "",
		Countdown:        20,
		BetRate:          "10/20",
		PlayerSize:       8,
		CardsOnTable:     []string{},
		PlayerListBySeat: make([]domain.PlayerDO, 0),
	}
}

func makeTableInfo() outbound.TableInfo {
	gi := outbound.TableInfo{
		TableID:      util.GenID(),
		TotalPot:     1000,
		Status:       "",
		Countdown:    20,
		BetRate:      "10/20",
		CardsOnTable: []string{"312", "412", "109"},
		DataType:     "table_info",
	}

	return gi
}
