package base

import (
	"errors"
	"time"

	"github.com/hudgit2019/leafboot/conf"

	"github.com/name5566/leaf/timer"
)

const (
	Tabletimer_checkbegin = iota
)

type TableTimer struct {
	t *timer.Timer
	d time.Duration
}
type GameTable struct {
	GameRoundNum   int64 //游戏对局编号
	SitdownPlayers int
	ReadyPlayers   int
	TimeGameBegin  time.Time
	TableStatus    int
	TablePlayers   []IPlayerNode
	WatchPlayers   []IPlayerNode
	Tabletimemap   map[int]*TableTimer //定时器
	gameLogic      IGameLogic
}

func (table *GameTable) Init(chairnum int, flogic IGameLogic) {
	table.TablePlayers = make([]IPlayerNode, chairnum)
	for j := 0; j < chairnum; j++ {
		table.TablePlayers[j] = nil
	}
	table.Tabletimemap = make(map[int]*TableTimer)
	table.gameLogic = flogic
}
func (table *GameTable) ResetTable() {
	for _, playerint := range table.TablePlayers {
		if playerint != nil {
			playerint.(*PlayerNode).Resetgamedata()
		}
	}
	for _, v := range table.Tabletimemap {
		if v.t != nil {
			v.t.Stop()
		}
	}
	table.ResetGameTable()
}
func (table *GameTable) ResetGameTable() {

}
func (table *GameTable) OnTimerCheckBegin() {
	if table.ReadyPlayers >= conf.RoomInfo.GameStartPlayer {
		table.gameLogic.Gamestart(table)
		table.Tabletimemap[Tabletimer_checkbegin] = nil
	} else {
		//continue checktimer
		t, _ := table.Tabletimemap[Tabletimer_checkbegin]
		table.SetTimer(Tabletimer_checkbegin, t.d, table.OnTimerCheckBegin)
	}
}
func (table *GameTable) KillTimer(timerid int) {
	if t, ok := table.Tabletimemap[timerid]; ok {
		if t.t != nil {
			t.t.Stop()
		}
		table.Tabletimemap[timerid] = nil
	}
}
func (table *GameTable) SetTimer(timerid int, d time.Duration, cb func()) error {
	if _, ok := table.Tabletimemap[timerid]; ok {
		return errors.New("timeid is already in use")
	}
	switch timerid {
	case Tabletimer_checkbegin:
		tatbletimer := &TableTimer{Skeleton.AfterFunc(d, cb), d}
		table.Tabletimemap[timerid] = tatbletimer
		break
	}
	return nil
}
func (table *GameTable) GameBegin(gamenum int64) {
	table.TableStatus = 1
	table.GameRoundNum = gamenum
	table.TimeGameBegin = time.Now()
	table.ResetTable()
	for _, playerint := range table.TablePlayers {
		if playerint != nil {
			playerint.(*ClientNode).Usergamestatus = PlayerstatuBeginGame
			playerint.(*ClientNode).GameBegin()
		}
	}
	for k, v := range table.Tabletimemap {
		if v.t != nil {
			v.t.Stop()
		}
		table.Tabletimemap[k] = nil
	}
}

func (table *GameTable) GameEnd() {
	table.ReadyPlayers = 0
	table.TableStatus = 0
	for k, v := range table.Tabletimemap {
		if v.t != nil {
			v.t.Stop()
		}
		table.Tabletimemap[k] = nil
	}
	for _, playerint := range table.TablePlayers {
		if playerint != nil {
			playerint.(*ClientNode).Usergamestatus = PlayerstatuHaveSitDown
			playerint.(*ClientNode).GameEnd()
		}
	}
}
