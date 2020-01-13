package base

import "github.com/hudgit2019/leafboot/msg"

const (
	PlayerstatuOffline     int32 = -1
	PlayerstatuInitial     int32 = 0
	PlayerstatuWaitAuthen  int32 = 1
	PlayerstatuWaitSitDown int32 = 2
	PlayerstatuHaveSitDown int32 = 3
	PlayerstatuWatching    int32 = 4
	PlayerstatuHandUp      int32 = 5
	PlayerstatuBeginGame   int32 = 6
)

type Userplaygamedata struct {
	Gametax        uint64                                 //税收
	Gamecoin       int64                                  //key:srcid 增量
	Goldbean       int32                                  //key:srcid增量
	Proplist       map[int32]map[int32]msg.UserPropChange //key:srcid道具列表增量
	Gameplaytime   int32                                  //
	Gameonlinetime int32                                  //
	Gamewintimes   int32                                  //
	Gamelosetimes  int32                                  //
	Gamecoinwin    int64                                  //
	Gamecoinlose   int64                                  //
}
type IPlayerNode interface {
	Initialize()
	IsProxyNode() bool
	IsMiddlePlatNode() bool
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
func (playerlist *PlayerNodeList) GetAllPlayers() map[int64]IPlayerNode {
	return playerlist.onlinelist
}
