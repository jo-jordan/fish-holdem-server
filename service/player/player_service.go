package player_service

import (
	"github.com/jo-jordan/fish-holdem-server/entity/domain"
	"github.com/jo-jordan/fish-holdem-server/misc/global"
)

func ActionCall() {}

func ActionCheck(player *domain.PlayerDO) {
	tableAny, ok := global.TableMap.Load(player.TableID)
	if !ok {
		return
	}
	table := tableAny.(domain.TableDO)

	if player.ID != table.CurActionPlayer.ID {
		return
	}

}

func ActionLooper() {

}
