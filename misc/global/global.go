package global

import "github.com/jo-jordan/fish-holdem-server/entity/domain"

var PlayerMap map[string]domain.PlayerDO
var TableMap map[int64]domain.TableDO

func Init() {
	PlayerMap = make(map[string]domain.PlayerDO)
	TableMap = make(map[int64]domain.TableDO)
}
