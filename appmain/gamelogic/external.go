package gamelogic

import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/conf"
	"github.com/hudgit2019/leafboot/db"
	"github.com/hudgit2019/leafboot/model"
	"github.com/hudgit2019/leafboot/msg"
	myredis "github.com/hudgit2019/leafboot/redis"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

type FactoryGameLogic struct {
	Curgamenum    int64 //当前牌局计数，初始化为时间戳
	GameTables    []base.ITable
	MapLevelTable map[int32][]base.ITable
}

func (f *FactoryGameLogic) keepAlive() {
	for {
		for _, v := range base.PlayerList.GetAllPlayers() {
			playerNode := v.(*base.ClientNode)
			//前端逻辑服器,5分钟更新一次token
			if f.IsFrontend() && time.Now().Unix()-playerNode.LastUpateTokenTime.Unix() >= 300 {
				log.Debug("player:%v updatetoken", playerNode.Usernodeinfo.Userid)
				db.UpdateATokenTTF(playerNode.Usernodeinfo.Userid)
				playerNode.LastUpateTokenTime = time.Now()
			}
		}
		time.Sleep(10 * time.Second)
	}
}

//listenEtcdConf 监听所有配置更新
func (f *FactoryGameLogic) listenEtcdConf() {
	for {
		select {
		case roomFlag := <-conf.ChanRoomFlag:
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
		case etcdconf := <-conf.ChanChildConf:
			{
				f.CallBackEtcdConf(etcdconf.Action, etcdconf.Key, etcdconf.Value)
			}
		case dbinfo := <-conf.ChanDataBase:
			{
				db.OpenDBGroup(dbinfo)
			}
		case redisinfo := <-conf.ChanRedisInfo:
			{
				myredis.OpenRedisGroup(redisinfo)
			}
		}

	}
}
func (f *FactoryGameLogic) Start(netgate *gate.Gate) error {
	base.NetGate = netgate
	f.Curgamenum = time.Now().Unix()
	f.CreateRoom()
	conf.StartEtcd()
	err := db.StartDB()
	if err != nil {
		log.Fatal("StartDB err:%v", err)
	}
	myredis.StartRedis()
	go f.keepAlive()
	go f.listenEtcdConf()
	return nil
}
func (f *FactoryGameLogic) CreateRoom() error {
	for i := 1; i <= int(conf.RoomInfo.MaxTableNum); i++ {
		table := f.OnCreateTable()
		table.Init(int32(i), conf.RoomInfo.MaxTableChair, f)
		f.GameTables = append(f.GameTables, table)
	}
	f.MapLevelTable = make(map[int32][]base.ITable)
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
					NodeName: conf.Server.NodeName,
					NodeID:   conf.Server.NodeID,
					IsGray:   conf.Server.IsGray,
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
		if playernode.Usergamestatus >= base.PlayerstatuWaitAuthen {
			f.handleloginout(player)
			f.OnTableLeave(playernode)
			f.CallBackLogOut(player)
		}
		//如果是被代理节点,推送到代理解除路由
		if playernode.IsProxyedNode() {
			leaveRsp := msg.LeaveSceneRes{
				Route: conf.Server.NodeName,
			}
			base.SendRspMsg(playernode, leaveRsp)
			proxyNode := playernode.Netagent.UserData().(*base.ProxyNode)
			delete(proxyNode.MapClient, playernode.ProxyClientID)
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
	_, okid := playernode.Useraccountdbw.Proplist[propinfo.Propid]
	if !okid {
		playernode.Useraccountdbw.Proplist[propinfo.Propid] = make(map[int32]msg.UserPropChange)
	}
	propidlist, okid := playernode.Useraccountdbw.Proplist[propinfo.Propid]
	log.Debug("SavePlayerProp userid:%v propinfo:%v writesource:%v", playernode.Usernodeinfo.Userid, propinfo, writesource)
	if okid {
		tmpProp, oktype := propidlist[propinfo.Proptype]
		if oktype {
			if propinfo.Propnum != 0 {
				tmpProp.Propnum += propinfo.Propnum
				playernode.Useraccountdbw.Proplist[propinfo.Propid][propinfo.Proptype] = tmpProp
			}
			if propinfo.Proptime != 0 {
				tmpProp.Proptime += propinfo.Proptime
				playernode.Useraccountdbw.Proplist[propinfo.Propid][propinfo.Proptype] = tmpProp
			}
		} else {
			playernode.Useraccountdbw.Proplist[propinfo.Propid][propinfo.Proptype] = propinfo
		}
	} else {
		// proptypelist := make(map[int32]msg.UserPropChange)
		playernode.Useraccountdbw.Proplist[propinfo.Propid][propinfo.Proptype] = propinfo
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
func (f *FactoryGameLogic) SavePlayerGameEnd(player base.IPlayerNode, datachanged base.Userplaygamedata) {
	playernode := player.(*base.ClientNode)
	log.Debug("SavePlayerGameEnd userid:%v datachanged:%v", playernode.Usernodeinfo.Userid, datachanged)
	if datachanged.Gamecoin > 0 {
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
func (f *FactoryGameLogic) ApplyRobot(robot model.ApplyRobotInfo) (base.IPlayerNode, error) {
	if player, ok := base.PlayerList.GetPlayer(robot.RobotID); ok {
		return player, errors.New(fmt.Sprintf("RobotID:%v is already in use", robot.RobotID))
	}
	player := f.OnCreatePlayer("").(*base.ClientNode)
	player.Initialize()
	player.Userisrobot = true
	player.Usernodeinfo.Userid = robot.RobotID
	player.Usernodeinfo.GameCoin = robot.GameCoin
	player.Usernodeinfo.GoldBean = robot.PrizeTicket
	player.Usernodeinfo.VipExp = robot.VipExp
	player.Usernodeinfo.NickName = robot.NickName
	player.Usernodeinfo.Gender = robot.Gender
	player.Usergamestatus = base.PlayerstatuWaitSitDown
	loginRes := msg.ApplyRobotReq{
		UserID:   robot.RobotID,
		NickName: robot.NickName,
		GameCoin: robot.GameCoin,
		TableID:  robot.TableID,
		ChairID:  robot.ChairID,
	}
	base.SendRspMsg(player, loginRes)
	return player, nil
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
func (f *FactoryGameLogic) AppMsgCallBackInit(map[string]base.MsgHandler)                 {}
func (f *FactoryGameLogic) CallBackEtcdConf(action string, key string, value string)      {}
func (f *FactoryGameLogic) OnDestroy() {
	//清除所有在线
	for _, playerint := range base.PlayerList.GetAllPlayers() {
		f.ClosePlayer(playerint)
	}
}

//IsFrontend 是否是前端服务器
func (f *FactoryGameLogic) IsFrontend() bool {
	return true
}
