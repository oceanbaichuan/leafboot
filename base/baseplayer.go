package base

import "github.com/hudgit2019/leafboot/msg"

const (
	PlayerstatuOffline     int = -1
	PlayerstatuInitial     int = 0
	PlayerstatuWaitAuthen  int = 1
	PlayerstatuWaitSitDown int = 2
	PlayerstatuHaveSitDown int = 3
	PlayerstatuWatching    int = 4
	PlayerstatuHandUp      int = 5
	PlayerstatuBeginGame   int = 6
)

type Userplaygamedata struct {
	Gametax        int64                              //税收
	Gamecoin       int64                              //增量
	Goldbean       int                                //增量
	Proplist       map[int]map[int]msg.UserPropChange //道具列表增量
	Gameplaytime   int                                //
	Gameonlinetime int                                //
	Gamewintimes   int                                //
	Gamelosetimes  int                                //
	Gamecoinwin    int64                              //
	Gamecoinlose   int64                              //
}
type IPlayerNode interface {
	Initialize()
	IsProxyNode() bool
	IsProxyedNode() bool
	Resetgamedata()
	GameBegin()
	GameEnd()
	HandleAutoGame()
}
type PlayerNodeList struct {
	onlinelist map[int64]IPlayerNode
}

func (playerlist *PlayerNodeList) Init() {
	playerlist.onlinelist = make(map[int64]IPlayerNode)
}
func (playerlist *PlayerNodeList) GetPlayer(userid int64) (IPlayerNode, bool) {
	v, ok := playerlist.onlinelist[userid]
	return v, ok
}
func (playerlist *PlayerNodeList) AddPlayer(userid int64, player IPlayerNode) {
	playerlist.onlinelist[userid] = player
}

func (playerlist *PlayerNodeList) DeletePlayer(userid int64) {
	delete(playerlist.onlinelist, userid)
}
func (playerlist *PlayerNodeList) GetOnlineNum() int {
	return len(playerlist.onlinelist)
}
