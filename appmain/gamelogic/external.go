package gamelogic

import (
	"errors"
	"time"
	"unsafe"

	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/conf"
	"github.com/hudgit2019/leafboot/db"
	"github.com/hudgit2019/leafboot/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

type FactoryGameLogic struct {
	Curgamenum int64 //当前牌局计数，初始化为时间戳
	GameTables []base.ITable
}

func (f *FactoryGameLogic) keepAlive() {
	for {
		time.Sleep(10 * time.Second)
	}
}

//updateFlag 更新房间状态到连接层
func (f *FactoryGameLogic) updateFlag(chanFlag chan bool) {
	for {
		select {
		case roomFlag := <-chanFlag:
			{
				for _, v := range mapConnServer {
					//注册本节点到connector
					flagNotice := msg.RequestData{
						Route: "Connector.GameFlag",
						ReqID: v.PlayerID,
						ReqData: &msg.GameFlagNotice{
							IsClosed: roomFlag,
						},
					}
					base.SendReqMsg(v, &flagNotice)
				}
			}
		}
	}
}
func (f *FactoryGameLogic) Start(netgate *gate.Gate) error {
	base.NetGate = netgate
	f.Curgamenum = time.Now().Unix()
	f.CreateRoom()
	conf.StartEtcd()
	db.StartDB()
	go f.keepAlive()
	go f.updateFlag(conf.ChanRoomFlag)
	return nil
}
func (f *FactoryGameLogic) CreateRoom() error {
	for i := 0; i < int(conf.RoomInfo.MaxTableNum); i++ {
		table := f.OnCreateTable()
		table.Init(conf.RoomInfo.MaxTableChair, f)
		f.GameTables = append(f.GameTables, table)
	}
	return nil
}
func (f *FactoryGameLogic) OnCreateTable() base.ITable {
	return &base.GameTable{}
}
func (f *FactoryGameLogic) GetTable(tableid int32) (base.ITable, error) {
	if tableid < 1 || tableid > conf.RoomInfo.MaxTableNum {
		return nil, errors.New("table index out of range")
	}
	return f.GameTables[tableid-1], nil
}
func (f *FactoryGameLogic) OnCreatePlayer(addr string) base.IPlayerNode {
	if _, ok := conf.MapConnServer[addr]; ok {
		return &base.ProxyNode{}
	} else {
		return f.CreateClientPlayer()
	}

}
func (f *FactoryGameLogic) CreateClientPlayer() base.IPlayerNode {
	return &base.ClientNode{}
}

func (f *FactoryGameLogic) OnPlayerConnect(player base.IPlayerNode) {
	if player.IsProxyNode() {
		playerNode := player.(*base.ProxyNode)
		if k, ok := conf.MapConnServer[playerNode.Netagent.RemoteAddr().String()]; ok {
			mapConnServer[k.LocalAddr] = playerNode
			//注册本节点到connector
			regReq := msg.RequestData{
				Route: "Connector.GameRegiste",
				ReqID: playerNode.PlayerID,
				ReqData: &msg.GameRegistReq{
					Addr:     conf.Server.WSAddr,
					NodeName: conf.RoomInfo.NodeName,
					NodeID:   conf.RoomInfo.NodeID,
				},
			}
			base.SendReqMsg(player, &regReq)
		}
		log.Debug("player: %s connected", playerNode.Netagent.RemoteAddr().String())
	} else {
		playerNode := player.(*base.ClientNode)
		log.Debug("player: %s connected", playerNode.Netagent.RemoteAddr().String())

	}
}

type iface struct {
	itype  uintptr
	ivalue uintptr
}

func (f *FactoryGameLogic) OnPlayerClose(player base.IPlayerNode) {
	//log.Debug("player: %v closed", player)
	if player == nil {
		return
	}
	if player.IsProxyNode() {
		playernode := player.(*base.ProxyNode)
		//被代理节点全部模拟掉线
		for key, v := range playernode.MapClient {
			delete(playernode.MapClient, key)
			f.OnPlayerClose(v)
		}
		delete(mapConnServer, playernode.Netagent.RemoteAddr().String())
		playernode.Netagent = nil
		return
	}
	playernode := player.(*base.ClientNode)
	if playernode.Usergamestatus < base.PlayerstatuBeginGame {
		log.Debug("playernode:%v 非游戏状态退出 Usergamestatus:%v", (*iface)(unsafe.Pointer(&playernode.Netagent)), playernode.Usergamestatus)
		if playernode.Usergamestatus > base.PlayerstatuWaitAuthen {
			f.handleloginout(player)
			f.OnTableLeave(playernode)
			f.CallBackLogOut(player)
		}
		//如果是被代理节点,推送到代理解除路由
		if playernode.IsProxyedNode() {
			leaveRsp := msg.LeaveScenceRes{
				Route: conf.RoomInfo.NodeName,
			}
			base.SendRspMsg(playernode, leaveRsp)
		}
	} else {
		log.Debug("playernode:%v 游戏状态退出 Usergamestatus:%v", playernode.Usernodeinfo.Userid, playernode.Usergamestatus)
		playernode.Useroffline = true
		f.CallBackOffline(player)
		player.HandleAutoGame()
	}
	playernode.Netagent = nil
}

func (f *FactoryGameLogic) ClosePlayer(player base.IPlayerNode) {
	if player.IsProxyedNode() { //如果是代理节点，立即释放
		f.OnPlayerClose(player)
	} else {
		playernode := player.(*base.PlayerNode)
		if playernode.Netagent != nil {
			playernode.Netagent.Close()
		}
	}
}

//GameStart 游戏对局开始处理，不要重写
func (f *FactoryGameLogic) Gamestart(table base.ITable) {
	f.Curgamenum++
	log.Debug("GameStart")
	table.GameBegin(f.Curgamenum)
	f.CallBackGameStart(table)
}

//SendRoomInfo 同步房间通用信息到c端，不要重写
func (f *FactoryGameLogic) SendRoomInfo(player base.IPlayerNode) {
	f.CallBackSendRoomInfo(player)
}

//SavePlayergamecoin 游戏结束调用，用户游戏过程中产生的金币变化更新
func (f *FactoryGameLogic) SavePlayerGameCoin(player base.IPlayerNode, gamecoin int64, writesource int32) {
	playernode := player.(*base.ClientNode)
	playernode.Useraccountdbw.Gamecoin += gamecoin
	f.WriteAttributionLog(&msg.AttributeChangelog{})
	playernode.Usernodeinfo.GameCoin += gamecoin
}

//SavePlayergoldbean 游戏结束调用，用户游戏过程中产生的金豆变化更新
func (f *FactoryGameLogic) SavePlayerGoldBean(player base.IPlayerNode, goldbean int32, writesource int32) {
	playernode := player.(*base.ClientNode)
	playernode.Useraccountdbw.Goldbean += goldbean
	f.WriteAttributionLog(&msg.AttributeChangelog{})
	playernode.Usernodeinfo.GoldBean += goldbean

}

//SavePlayerprop 游戏结束调用，用户游戏过程中产生的道具变化更新
func (f *FactoryGameLogic) SavePlayerProp(player base.IPlayerNode, propinfo msg.UserPropChange, writesource int32) {
	playernode := player.(*base.ClientNode)
	propidlist, okid := playernode.Useraccountdbw.Proplist[propinfo.Propid]
	log.Debug("SavePlayerProp userid:%v propinfo:%v writesource:%v", playernode.Usernodeinfo.Userid, propinfo, writesource)
	if okid {
		propinfo, oktype := propidlist[propinfo.Proptype]
		if oktype {
			if propinfo.Propnum != 0 {
				propinfo.Propnum += propinfo.Propnum

			}
			if propinfo.Proptime != 0 {
				propinfo.Proptime += propinfo.Proptime
			}
		} else {
			propidlist[propinfo.Proptype] = propinfo
		}
	} else {
		proptypelist := make(map[int32]msg.UserPropChange)
		proptypelist[propinfo.Proptype] = propinfo
		playernode.Useraccountdbw.Proplist[propinfo.Propid] = proptypelist
	}
	bExists := false
	for _, curpropinfo := range playernode.Usernodeinfo.Proplist {
		if curpropinfo.Propid == propinfo.Propid &&
			curpropinfo.Proptype == propinfo.Proptype {
			if propinfo.Propnum != 0 {
				f.WriteAttributionLog(&msg.AttributeChangelog{})
				curpropinfo.Propnum += propinfo.Propnum
			} else if propinfo.Proptime != 0 {
				f.WriteAttributionLog(&msg.AttributeChangelog{})
				curpropinfo.Proptime = curpropinfo.Proptime + propinfo.Proptime
			}
			bExists = true
			break
		}
	}
	if !bExists {
		tempprop := msg.PropInfo{
			Propid:   propinfo.Propid,
			Proptype: propinfo.Proptype,
			Propnum:  propinfo.Propnum,
			Proptime: propinfo.Proptime,
		}
		playernode.Usernodeinfo.Proplist = append(playernode.Usernodeinfo.Proplist, tempprop)
		if propinfo.Propnum != 0 {
			f.WriteAttributionLog(&msg.AttributeChangelog{})
		} else if propinfo.Proptime != 0 {
			f.WriteAttributionLog(&msg.AttributeChangelog{})
		}
	}
}

//SavePlayergameend 游戏结束调用，用户游戏过程中产生的金豆变化更新
func (f *FactoryGameLogic) SavePlayerGameEnd(player base.IPlayerNode, datachanged base.Userplaygamedata, writesource int32) {
	playernode := player.(*base.ClientNode)
	log.Debug("SavePlayerGameEnd userid:%v datachanged:%v", playernode.Usernodeinfo.Userid, datachanged)
	if datachanged.Gamecoin != 0 {
		f.WriteAttributionLog(&msg.AttributeChangelog{})
	}
	if datachanged.Goldbean != 0 {
		f.WriteAttributionLog(&msg.AttributeChangelog{})
	}
	//逻辑值更新
	playernode.Usernodeinfo.GameCoin += datachanged.Gamecoin + datachanged.Gametax
	playernode.Usernodeinfo.GoldBean += datachanged.Goldbean
	playernode.Usernodeinfo.GameCoinWin += datachanged.Gamecoin

	//增量值更新
	playernode.Useraccountdbw.Gamecoin += datachanged.Gamecoin + datachanged.Gametax
	playernode.Useraccountdbw.Goldbean += datachanged.Goldbean
	playernode.Useraccountdbw.Gamecoinwin += datachanged.Gamecoin
	if playernode.Useraccountdbw.Gamecoin > 0 {
		playernode.Useraccountdbw.Gamewintimes += 1
		playernode.Usernodeinfo.GameWinTimes += 1
	} else {
		playernode.Useraccountdbw.Gamecoinlose += 1
		playernode.Usernodeinfo.GameCoinLose += 1
	}
	playernode.Useraccountdbw.Gametax += datachanged.Gametax
	tableint, _ := f.GetTable(playernode.Usertableid)
	if tableint != nil {
		tableitem := tableint.(*base.GameTable)
		playtime := (int32)(time.Now().Unix() - tableitem.TimeGameBegin.Unix())
		playernode.Useraccountdbw.Gameplaytime += playtime
		playernode.Usernodeinfo.GamePlayTime += playtime
	}
}

//SaveTableGameEnd 结束桌子所有对局状态，在SavePlayerGameEnd、WriteTableRoundLog完成后调用
func (f *FactoryGameLogic) SaveTableGameEnd(table base.ITable) {
	gametable := table.(*base.GameTable)
	//对局日志记录
	for _, playerint := range gametable.TablePlayers {
		if playerint != nil {
			//断线用户，踢出桌子
			playerNode := playerint.(*base.ClientNode)
			if playerNode.Useroffline {
				playerNode.Usergamestatus = base.PlayerstatuHaveSitDown
				f.OnPlayerClose(playerint)
			} else {
				playerNode.Usergamestatus = base.PlayerstatuHaveSitDown
			}
		}
	}
	gametable.GameEnd()
}
func (f *FactoryGameLogic) HandleAutoGame(player base.IPlayerNode) {
	f.AutoPlay(player)
}

func (f *FactoryGameLogic) CallBackSendRoomInfo(player base.IPlayerNode)                  {}
func (f *FactoryGameLogic) CallBackSitDown(player base.IPlayerNode, table base.ITable)    {}
func (f *FactoryGameLogic) CallLoginSuccess(player base.IPlayerNode)                      {}
func (f *FactoryGameLogic) CallBackLeaveTable(player base.IPlayerNode, table base.ITable) {}
func (f *FactoryGameLogic) CallBackOffline(player base.IPlayerNode)                       {}
func (f *FactoryGameLogic) CallBackLogOut(player base.IPlayerNode)                        {}
func (f *FactoryGameLogic) CallBackHandUp(player base.IPlayerNode, table base.ITable)     {}
func (f *FactoryGameLogic) CallBackGameStart(table base.ITable)                           {}
func (f *FactoryGameLogic) CallBackLoginAgain(player base.IPlayerNode)                    {}
func (f *FactoryGameLogic) AutoPlay(player base.IPlayerNode)                              {}
func (f *FactoryGameLogic) AppMsgCallBackInit(*map[string]base.MsgHandler)                {}
func (f *FactoryGameLogic) OnDestroy() {
}
