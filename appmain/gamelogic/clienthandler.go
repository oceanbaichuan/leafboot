package gamelogic

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/conf"
	"github.com/hudgit2019/leafboot/db"
	"github.com/hudgit2019/leafboot/msg"

	"github.com/goinggo/mapstructure"
	"github.com/name5566/leaf/log"
)

func (f *FactoryGameLogic) handleLoginGameReq(args []interface{}) {
	clientReq := args[0].(*msg.RequestData)
	// 消息的发送者
	playerInt := args[1].(base.IPlayerNode)
	req := msg.EnterScenceReq{}
	err := mapstructure.Decode(clientReq.ReqData, &req)
	if err != nil {
		base.SendFailMsgWithID(playerInt, clientReq.ReqID, 405,
			fmt.Sprintf("invalid param for method:%s", clientReq.Route), nil)
		log.Error("invalid param for method:%s", clientReq.Route)
		return
	}
	log.Debug("login %v %v", req.UserID, req.Passwd)
	routeParam := strings.Split(clientReq.Route, ".")
	var tmpplayer *base.ClientNode = nil
	if playerInt.IsProxyNode() {
		proxyNode := playerInt.(*base.ProxyNode)
		proxyplayer, ok := proxyNode.MapClient[clientReq.ReqID]
		if !ok {
			newplayer := f.OnCreatePlayer("").(*base.ClientNode)
			newplayer.Initialize()
			newplayer.ProxyClientID = clientReq.ReqID
			newplayer.Netagent = proxyNode.Netagent
			tmpplayer = newplayer
			proxyNode.MapClient[clientReq.ReqID] = newplayer
		} else {
			tmpplayer = proxyplayer
		}
	} else {
		tmpplayer = playerInt.(*base.ClientNode)
	}
	if len(routeParam) < 2 {
		base.SendFailMsg(tmpplayer, 404, fmt.Sprintf("route:%s undefined", clientReq.Route), nil)
		return
	}
	//登录过，不再做登录加载
	if tmpplayer.Usergamestatus > base.PlayerstatuWaitAuthen {
		return
	}
	var succplayer *base.ClientNode = tmpplayer
	//正常登录
	//succplayer := a.UserData().(*base.PlayerNode)
	succplayer.Userlogindeviceid = req.MacID
	//succplayer.Userlogindevicetype = req.Logindevicetype
	//succplayer.Userloginchannellong = req.Loginlongchannel
	//succplayer.Userloginchannelshort = req.Loginshortchannel
	succplayer.Userloginfrom = req.AppID
	succplayer.LoginIP = req.LoginIP
	succplayer.Usergamestatus = base.PlayerstatuWaitAuthen
	// succplayer.Userloginsiteid = req.Loginsiteid
	//token验证,告诉用户登录失败
	// loginres := msg.Logingamebaseres{
	// 	Errcode: msg.LoginErr_tokenerr,
	// }
	// base.Sendmsg(succplayer, loginres)
	// a.Close()
	dbloginres := msg.LoginRes{}
	dbloginreq := msg.Accountrdatareq{
		Account:     req.Account,
		Passwd:      req.Passwd,
		Token:       req.Token,
		Userid:      req.UserID,
		Gameid:      conf.RoomInfo.GameID,
		Serverid:    conf.RoomInfo.NodeID,
		Loginip:     strings.Split(req.LoginIP, ":")[0],
		AppID:       req.AppID,
		MacID:       req.MacID,
		ServerLevel: conf.RoomInfo.RoomLevel,
		Proplist:    conf.RoomInfo.PropIdList,
	}
	err = db.SelectAccount(&dbloginreq, &dbloginres)
	if err == nil {
		log.Debug("SelectAccount err:%v", err)
		db.SelectUserInfo(&dbloginreq, &dbloginres)
		log.Debug("SelectUserInfo err:%v", err)
		db.SelectGameData(&dbloginreq, &dbloginres)
		log.Debug("SelectGameData err:%v", err)
	}
	//登录成功后，添加节点
	f.handleLoginRes([]interface{}{&dbloginres, succplayer})
}
func (f *FactoryGameLogic) handleLoginAgain(player base.IPlayerNode) {
	succplayer := player.(*base.ClientNode)
	loginagainres := msg.LoginAgainRes{
		UserID:     succplayer.Usernodeinfo.Userid,
		ErrCode:    succplayer.Usernodeinfo.Errcode,
		Account:    succplayer.Usernodeinfo.Account,
		Token:      succplayer.Usernodeinfo.Token,
		AcIndex:    succplayer.Usernodeinfo.AcIndex,
		AcType:     succplayer.Usernodeinfo.AcType,
		AppID:      succplayer.Usernodeinfo.AppID,
		IsValid:    succplayer.Usernodeinfo.IsValid,
		SID:        succplayer.Usernodeinfo.SiteID,
		NickName:   succplayer.Usernodeinfo.NickName,
		GameCoin:   succplayer.Usernodeinfo.GameCoin,
		PTicket:    succplayer.Usernodeinfo.GoldBean,
		AllPTicket: succplayer.Usernodeinfo.AllGoldBean,
		RegChan:    succplayer.Usernodeinfo.RegChan,
		RegTime:    succplayer.Usernodeinfo.RegTime,
		RegSID:     succplayer.Usernodeinfo.RegSiteID,
		Gender:     succplayer.Usernodeinfo.Gender,
		HeadID:     succplayer.Usernodeinfo.HeadID,
		THUrl:      succplayer.Usernodeinfo.ThirdHeadUrl,
		PhNum:      succplayer.Usernodeinfo.Phonebinded,
		GameID:     succplayer.Usernodeinfo.OtherGameID,
		ServerID:   succplayer.Usernodeinfo.OtherRoomID,
		PropList:   succplayer.Usernodeinfo.Proplist,
	}
	//入座状态
	if succplayer.Usergamestatus >= base.PlayerstatuHaveSitDown {
		gametable := f.GameTables[succplayer.Usertableid].(*base.GameTable)
		for _, player := range gametable.TablePlayers {
			playernode := player.(*base.ClientNode)
			tableplayer := msg.TablePlayer{}
			tableplayer.Userid = playernode.Usernodeinfo.Userid
			tableplayer.Siteid = playernode.Usernodeinfo.RegSiteID
			tableplayer.Tableid = playernode.Usertableid
			tableplayer.Chairid = playernode.Userchairid
			tableplayer.Gamecoin = playernode.Usernodeinfo.GameCoin
			tableplayer.Bankcoin = playernode.Usernodeinfo.BankCoin
			tableplayer.Gamelosetimes = playernode.Usernodeinfo.GameLoseTimes
			tableplayer.Gamestatus = playernode.Usergamestatus
			tableplayer.Gamewintimes = playernode.Usernodeinfo.GameWinTimes
			tableplayer.Gender = playernode.Usernodeinfo.Gender
			tableplayer.Goldbean = playernode.Usernodeinfo.GoldBean
			tableplayer.Nickname = playernode.Usernodeinfo.NickName
			tableplayer.Sysheadid = playernode.Usernodeinfo.HeadID
			tableplayer.Thirdheadurl = playernode.Usernodeinfo.ThirdHeadUrl
			tableplayer.Vipexp = playernode.Usernodeinfo.VipExp
			loginagainres.Tableplayers = append(loginagainres.Tableplayers, tableplayer)
		}
	}
	base.SendRspMsg(succplayer, loginagainres)
	f.CallBackLoginAgain(player)
}
func (f *FactoryGameLogic) handleLoginRes(args []interface{}) {
	succplayer := args[1].(*base.ClientNode)
	dbloginres := args[0].(*msg.LoginRes)
	succplayer.Usernodeinfo = *dbloginres
	clientLoginRsp := msg.ClientUserRes{
		UserID:     dbloginres.Userid,
		ErrCode:    dbloginres.Errcode,
		Account:    dbloginres.Account,
		Token:      dbloginres.Token,
		AcIndex:    dbloginres.AcIndex,
		AcType:     dbloginres.AcType,
		AppID:      dbloginres.AppID,
		IsValid:    dbloginres.IsValid,
		SID:        dbloginres.SiteID,
		NickName:   dbloginres.NickName,
		GameCoin:   dbloginres.GameCoin,
		PTicket:    dbloginres.GoldBean,
		AllPTicket: dbloginres.AllGoldBean,
		RegChan:    dbloginres.RegChan,
		RegTime:    dbloginres.RegTime,
		RegSID:     dbloginres.RegSiteID,
		Gender:     dbloginres.Gender,
		HeadID:     dbloginres.HeadID,
		THUrl:      dbloginres.ThirdHeadUrl,
		PhNum:      dbloginres.Phonebinded,
		GameID:     dbloginres.OtherGameID,
		ServerID:   dbloginres.OtherRoomID,
		PropList:   dbloginres.Proplist,
	}
	if dbloginres.IsValid == 0 {
		clientLoginRsp.ErrCode = msg.LoginErr_accountforbidden
		base.SendRspMsg(succplayer, clientLoginRsp)
		log.Error("userid:%v IsInValid", dbloginres.Userid)
		f.ClosePlayer(succplayer)
		return
	} else if dbloginres.OtherGameID > 0 && (dbloginres.OtherGameID != conf.RoomInfo.GameID ||
		dbloginres.OtherRoomID != conf.RoomInfo.NodeID) {
		log.Error("userid:%v in other room %v %v", dbloginres.Userid, dbloginres.OtherGameID, dbloginres.OtherRoomID)
		//告诉用户登录失败
		clientLoginRsp.ErrCode = msg.LoginErr_inotherroom
		base.SendRspMsg(succplayer, clientLoginRsp)
		f.ClosePlayer(succplayer)
		return
	} else if dbloginres.OtherGameID > 0 &&
		dbloginres.OtherGameID == conf.RoomInfo.GameID &&
		dbloginres.OtherRoomID == conf.RoomInfo.NodeID {
		//检测是否重复登录,踢掉老节点
		oldplayerintf, ok := base.PlayerList.GetPlayer(dbloginres.Userid)
		if ok {
			//登录验证
			oldplayer := oldplayerintf.(*base.ClientNode)
			//copy oldplayer to newplayer
			if oldplayer.Usergamestatus >= base.PlayerstatuBeginGame {
				//清空老节点的游戏数据
				oldplayer.Netagent.SetUserData(nil)
				f.ClosePlayer(oldplayerintf)
				oldplayer.Userlogindeviceid = succplayer.Userlogindeviceid
				oldplayer.Userlogindevicetype = succplayer.Userlogindevicetype
				oldplayer.Userloginchannellong = succplayer.Userloginchannellong
				oldplayer.Userloginchannelshort = succplayer.Userloginchannelshort
				oldplayer.Userloginfrom = succplayer.Userloginfrom
				oldplayer.Userloginsiteid = succplayer.Userloginsiteid
				oldplayer.LoginIP = succplayer.LoginIP
				oldplayer.Useroffline = false
				oldplayer.PlayerID = succplayer.PlayerID
				oldplayer.ProxyClientID = succplayer.ProxyClientID
				oldplayer.Netagent = succplayer.Netagent
				//f.CopyPlayerNode(oldplayer, succplayer)
				//a.SetUserData(oldplayer)
				//断线重连同步数据到前端
				f.handleLoginAgain(oldplayer)
				return
			} else if oldplayer.Usergamestatus >= base.PlayerstatuHaveSitDown &&
				oldplayer.Usergamestatus < base.PlayerstatuBeginGame {
				f.OnTableLeave(oldplayer)
				oldplayer.Usergamestatus = base.PlayerstatuWaitAuthen
				f.ClosePlayer(oldplayerintf)
			} else {
				oldplayer.Usergamestatus = base.PlayerstatuWaitAuthen
				f.ClosePlayer(oldplayerintf)
			}
			if oldplayer.Usergamestatus > base.PlayerstatuWaitAuthen &&
				oldplayer.Usergamestatus < base.PlayerstatuBeginGame {
				//将老节点属性的增量值同步到新节点
				succplayer.SyncOldPlayerData(oldplayer)
			}
		}
	} else if conf.RoomInfo.RoomCoinDownLimit != 0 && dbloginres.GameCoin < conf.RoomInfo.RoomCoinDownLimit {
		clientLoginRsp.ErrCode = msg.LoginErr_gamecoin_notenough
		base.SendRspMsg(succplayer, clientLoginRsp)
		log.Error("userid:%v money:%v not enough.limit:%v", dbloginres.Userid, dbloginres.GameCoin, conf.RoomInfo.RoomCoinDownLimit)
		f.ClosePlayer(succplayer)
		return
	} else if conf.RoomInfo.RoomCoinUpLimit != 0 && dbloginres.GameCoin > conf.RoomInfo.RoomCoinUpLimit {
		clientLoginRsp.ErrCode = msg.LoginErr_gamecoin_toomuch
		base.SendRspMsg(succplayer, clientLoginRsp)
		log.Error("userid:%v money:%v too much.limit:%v", dbloginres.Userid, dbloginres.GameCoin, conf.RoomInfo.RoomCoinUpLimit)
		f.ClosePlayer(succplayer)
		return
	}
	base.SendRspMsg(succplayer, succplayer.Usernodeinfo)
	base.PlayerList.AddPlayer(dbloginres.Userid, succplayer)
	//记录登录完成时间点
	succplayer.Userlastopertime = time.Now()
	succplayer.Userlogintime = time.Now()
	succplayer.Usergamestatus = base.PlayerstatuWaitSitDown
	f.CallLoginSuccess(succplayer)
	f.SendRoomInfo(succplayer)
	f.WriteLoginRoomLog(&msg.Logingamelog{})
}

func (f *FactoryGameLogic) handleSitdownReq(args []interface{}) {
	clientReq := args[0].(*msg.RequestData)
	// 消息的发送者
	playerInt := args[1].(base.IPlayerNode)
	req := msg.SitDownRes{}
	mapstructure.Decode(clientReq.ReqData, &req)
	log.Debug("sitdown %v", req)
	routeParam := strings.Split(clientReq.Route, ".")
	var tmpplayer *base.ClientNode = nil
	if playerInt.IsProxyNode() {
		proxyNode := playerInt.(*base.ProxyNode)
		proxyplayer, ok := proxyNode.MapClient[clientReq.ReqID]
		if !ok {
			sitres := msg.SitDownRes{}
			sitres.Errcode = msg.SitErr_NoLogin
			base.SendRspMsgWithID(proxyNode, clientReq.ReqID, sitres)
			return
		} else {
			tmpplayer = proxyplayer
		}
	} else {
		tmpplayer = playerInt.(*base.ClientNode)
	}
	if len(routeParam) < 2 {
		base.SendFailMsg(tmpplayer, 404, fmt.Sprintf("route:%s undefined", clientReq.Route), nil)
		return
	}
	player := tmpplayer
	sitres := msg.SitDownRes{}
	log.Debug("handleSitdownReq playerstatu:%v req:%v", player.Usergamestatus, req)
	if player.Usergamestatus != base.PlayerstatuWaitSitDown {
		sitres.Errcode = msg.SitErr_NoLogin
		base.SendRspMsg(player, sitres)
		return
	}
	if player.Usergamestatus > base.PlayerstatuWaitSitDown {
		sitres.Errcode = msg.SitErr_HaveSit
		base.SendRspMsg(player, sitres)
		return
	}
	f.ArrangePlayerSitDownReq(player, req.Tableid, req.Chairid)
}
func (f *FactoryGameLogic) ArrangePlayerSitDownReq(player *base.ClientNode, tableid int32, chairid int32) {
	var fittable base.ITable
	sitres := msg.SitDownRes{}
	if tableid > 0 && tableid <= conf.RoomInfo.MaxTableNum {
		gametable, rettableid, retchairid, err := f.FixSearchTable(player, tableid, chairid)
		sitres.Tableid = rettableid
		sitres.Chairid = retchairid
		if err != nil {
			base.SendRspMsg(player, sitres)
			log.Error("userid:%v fixtable:%v:%v sitdown err:%v", player.Usernodeinfo.Userid, tableid, chairid, err)
			return
		}
		fittable = gametable

	} else if tableid > conf.RoomInfo.MaxTableNum || tableid < 1 {
		gametable, rttableid, rtchairid, err := f.AutoSearchTable(player)
		sitres.Tableid = rttableid
		sitres.Chairid = rtchairid
		if err != nil {
			base.SendRspMsg(player, sitres)
			log.Error("userid:%v auto sitdown err:%v", player.Usernodeinfo.Userid, err)
			return
		}
		fittable = gametable
	}
	tableitem := fittable.(*base.GameTable)
	tableitem.SitdownPlayers++
	for _, playerint := range tableitem.TablePlayers {
		if playerint != nil {
			playernode := playerint.(*base.ClientNode)
			tableplayer := msg.TablePlayer{}
			tableplayer.Userid = playernode.Usernodeinfo.Userid
			tableplayer.Siteid = playernode.Usernodeinfo.RegSiteID
			tableplayer.Tableid = playernode.Usertableid
			tableplayer.Chairid = playernode.Userchairid
			tableplayer.Gamecoin = playernode.Usernodeinfo.GameCoin
			tableplayer.Bankcoin = playernode.Usernodeinfo.BankCoin
			tableplayer.Gamelosetimes = playernode.Usernodeinfo.GameLoseTimes
			tableplayer.Gamestatus = playernode.Usergamestatus
			tableplayer.Gamewintimes = playernode.Usernodeinfo.GameWinTimes
			tableplayer.Gender = playernode.Usernodeinfo.Gender
			tableplayer.Goldbean = playernode.Usernodeinfo.GoldBean
			tableplayer.Nickname = playernode.Usernodeinfo.NickName
			tableplayer.Sysheadid = playernode.Usernodeinfo.HeadID
			tableplayer.Thirdheadurl = playernode.Usernodeinfo.ThirdHeadUrl
			tableplayer.Vipexp = playernode.Usernodeinfo.VipExp
			sitres.Players = append(sitres.Players, tableplayer)
		}
	}
	player.Userchairid = sitres.Chairid
	player.Usertableid = sitres.Tableid
	player.Usergamestatus = base.PlayerstatuHaveSitDown
	base.SendRspMsg(player, sitres)
	//通知其他人，玩家入座
	playerjoinres := msg.TableJoinPlayer{}
	playerjoinres.Player.Userid = player.Usernodeinfo.Userid
	playerjoinres.Player.Siteid = player.Usernodeinfo.RegSiteID
	playerjoinres.Player.Tableid = player.Usertableid
	playerjoinres.Player.Chairid = player.Userchairid
	playerjoinres.Player.Gamecoin = player.Usernodeinfo.GameCoin
	playerjoinres.Player.Bankcoin = player.Usernodeinfo.BankCoin
	playerjoinres.Player.Gamelosetimes = player.Usernodeinfo.GameLoseTimes
	playerjoinres.Player.Gamestatus = player.Usergamestatus
	playerjoinres.Player.Gamewintimes = player.Usernodeinfo.GameWinTimes
	playerjoinres.Player.Gender = player.Usernodeinfo.Gender
	playerjoinres.Player.Goldbean = player.Usernodeinfo.GoldBean
	playerjoinres.Player.Nickname = player.Usernodeinfo.NickName
	playerjoinres.Player.Sysheadid = player.Usernodeinfo.HeadID
	playerjoinres.Player.Thirdheadurl = player.Usernodeinfo.ThirdHeadUrl
	playerjoinres.Player.Vipexp = player.Usernodeinfo.VipExp
	for _, playerint := range tableitem.TablePlayers {
		if playerint != nil && playerint.(*base.ClientNode).Usernodeinfo.Userid != player.Usernodeinfo.Userid {
			base.SendRspMsg(playerint, playerjoinres)
		}
	}

	log.Debug("userid:%v sitdown success:%v:%v", player.Usernodeinfo.Userid, sitres.Tableid, sitres.Chairid)
	//记录入座时间点
	player.Userlastopertime = time.Now()
	player.SetTimer(base.Playertimer_checkhandup, time.Duration(conf.RoomInfo.SitNoHandUpCheckTime)*time.Second)
	if conf.RoomInfo.GameStartPlayer < conf.RoomInfo.MaxTableChair {
		tableitem.SetTimer(base.Tabletimer_checkbegin, time.Duration(conf.RoomInfo.GameStartCheckTime)*time.Second, tableitem.OnTimerCheckBegin)
	}
	f.CallBackSitDown(player, fittable)
}
func (f *FactoryGameLogic) AutoSearchTable(playernode *base.ClientNode) (base.ITable, int32, int32, error) {
	var oktable base.ITable
	var err error
	tablechair := -1
	tableid := -1
	for i, tableint := range f.GameTables {
		table := tableint.(*base.GameTable)
		if table.TableStatus == 1 {
			continue
		}
		for j, playerint := range table.TablePlayers {
			if playerint == nil {
				oktable = table
				tableid = i + 1
				tablechair = j + 1
				table.TablePlayers[j] = playernode
				break
			}
		}
		if oktable != nil {
			break
		}
	}
	if oktable == nil {
		err = errors.New("no fit table found")
	}
	return oktable, int32(tableid), int32(tablechair), err
}

func (f *FactoryGameLogic) FixSearchTable(playernode *base.ClientNode, tableid int32, chairid int32) (base.ITable, int32, int32, error) {
	var err error
	var oktable base.ITable
	gametableint, _ := f.GetTable(tableid)
	fittable := gametableint.(*base.GameTable)
	if conf.RoomInfo.CanWatchGame == 0 && fittable.TableStatus == 1 {
		return nil, -1, -1, errors.New("game have started forbidden to watch!")
	}
	for j, playerint := range fittable.TablePlayers {
		//log.Debug("j:%v playerint:%v", j, playerint)
		if int32(j) == chairid-1 && playerint == nil {
			fittable.TablePlayers[j] = playernode
			oktable = fittable
			chairid = int32(j) + 1
			break
		} else if chairid < 1 && playerint == nil {
			fittable.TablePlayers[j] = playernode
			oktable = fittable
			chairid = int32(j) + 1
			break
		}
	}
	if oktable == nil {

		return f.AutoSearchTable(playernode)
	}
	return oktable, tableid, chairid, err
}
func (f *FactoryGameLogic) handleHandUpReq(args []interface{}) {
	clientReq := args[0].(*msg.RequestData)
	// 消息的发送者
	playerInt := args[1].(base.IPlayerNode)
	handres := msg.HandUpRes{}
	//mapstructure.Decode(clientReq.ReqData, &req)
	routeParam := strings.Split(clientReq.Route, ".")
	var player *base.ClientNode = nil
	if playerInt.IsProxyNode() {
		proxyNode := playerInt.(*base.ProxyNode)
		proxyplayer, ok := proxyNode.MapClient[clientReq.ReqID]
		if !ok {
			handres.Errcode = msg.Handuperr_no_sitdown
			base.SendRspMsgWithID(proxyNode, clientReq.ReqID, handres)
			return
		} else {
			player = proxyplayer
		}
	} else {
		player = playerInt.(*base.ClientNode)
	}
	log.Debug("handup userid:%d", player.Usernodeinfo.Userid)
	if len(routeParam) < 2 {
		base.SendFailMsg(player, 404, fmt.Sprintf("route:%s undefined", clientReq.Route), nil)
		return
	}

	if player.Usergamestatus < base.PlayerstatuHaveSitDown {
		handres.Errcode = msg.Handuperr_no_sitdown //not sitdown
		base.SendRspMsg(player, handres)
		return
	}
	if player.Usergamestatus >= base.PlayerstatuHandUp {
		handres.Errcode = msg.Handuperr_already_hand //already hand
		base.SendRspMsg(player, handres)
		return
	}
	if conf.RoomInfo.RoomCoinDownLimit != 0 && player.Usernodeinfo.GameCoin < conf.RoomInfo.RoomCoinDownLimit {
		handres.Errcode = msg.Handuperr_coin_notenough //coin not enough
		base.SendRspMsg(player, handres)
		base.ClosePlayer(player)
		return
	} else if conf.RoomInfo.RoomCoinUpLimit != 0 && player.Usernodeinfo.GameCoin > conf.RoomInfo.RoomCoinUpLimit {
		handres.Errcode = msg.Handuperr_coin_toomuch //coin too much
		base.SendRspMsg(player, handres)
		base.ClosePlayer(player)
		return
	} else if conf.RoomInfo.ServerStatus == 0 {
		handres.Errcode = msg.Handuperr_room_closed //room closed
		base.SendRspMsg(player, handres)
		base.ClosePlayer(player)
		return
	}
	player.KillTimer(base.Playertimer_checkhandup)
	base.SendRspMsg(player, handres)
	handnotice := msg.HandUpNotice{}
	handnotice.Tableid = player.Usertableid
	handnotice.Chairid = player.Userchairid
	gametableint, _ := f.GetTable(player.Usertableid)
	gametable := gametableint.(*base.GameTable)
	for _, playerint := range gametable.TablePlayers {
		//log.Debug("handleHandUpReq player:%v playerint:%v", player, playerint)
		if playerint != nil && playerint.(*base.ClientNode).Usernodeinfo.Userid != player.Usernodeinfo.Userid {
			base.SendRspMsg(playerint, handnotice)
		}
	}
	//记录准备时间点
	player.Userlastopertime = time.Now()
	f.CallBackHandUp(player, gametable)
	gametable.ReadyPlayers++
	if gametable.ReadyPlayers == conf.RoomInfo.GameStartPlayer &&
		gametable.TableStatus == 0 {
		f.Gamestart(gametable)
	}
}
func (f *FactoryGameLogic) handleLeaveTableReq(args []interface{}) {
	playerInt := args[1].(base.IPlayerNode)
	clientReq := args[0].(*msg.RequestData)
	leaveres := msg.TablePlayerLeaveRes{}
	req := msg.TablePlayerLeaveReq{}
	mapstructure.Decode(clientReq.ReqData, &req)
	routeParam := strings.Split(clientReq.Route, ".")
	var player *base.ClientNode = nil
	if playerInt.IsProxyNode() {
		proxyNode := playerInt.(*base.ProxyNode)
		proxyplayer, ok := proxyNode.MapClient[clientReq.ReqID]
		if !ok {
			leaveres.Errcode = 1
			base.SendRspMsgWithID(proxyNode, clientReq.ReqID, leaveres)
			return
		} else {
			player = proxyplayer
		}
	} else {
		player = playerInt.(*base.ClientNode)
	}
	log.Debug("leavetable userid:%d", player.Usernodeinfo.Userid)
	if len(routeParam) < 2 {
		base.SendFailMsg(player, 404, fmt.Sprintf("route:%s undefined", clientReq.Route), nil)
		return
	}
	if player.Usergamestatus < base.PlayerstatuHaveSitDown || player.Usergamestatus > base.PlayerstatuHandUp {
		leaveres := msg.TablePlayerLeaveRes{}
		leaveres.Leavetype = req.Leavetype
		leaveres.Errcode = 1
		base.SendRspMsg(player, leaveres)
		return
	}
	f.OnTableLeave(player)
	if req.Leavetype == msg.TableLeave_ChangeTable {
		f.ArrangePlayerSitDownReq(player, 0, 0)
	}
}
func (f *FactoryGameLogic) OnTableLeave(player *base.ClientNode) {
	if player.Usergamestatus < base.PlayerstatuHaveSitDown ||
		player.Usergamestatus > base.PlayerstatuHandUp {
		return
	}
	gametableint, _ := f.GetTable(player.Usertableid)
	gametable := gametableint.(*base.GameTable)
	gametable.TablePlayers[player.Userchairid-1] = nil
	gametable.SitdownPlayers--
	if player.Usergamestatus == base.PlayerstatuHandUp {
		gametable.ReadyPlayers--
	}
	leavenotice := msg.TablePlayerLeaveNotice{}
	leavenotice.Userid = player.Usernodeinfo.Userid
	leavenotice.Tableid = player.Usertableid
	leavenotice.Chairid = player.Userchairid
	for _, playerint := range gametable.TablePlayers {
		if playerint != nil && playerint.(*base.ClientNode).Usernodeinfo.Userid != player.Usernodeinfo.Userid {
			base.SendRspMsg(playerint, leavenotice)
		}
	}
	leaveres := msg.TablePlayerLeaveRes{}
	leaveres.Leavetype = 1
	leaveres.Errcode = 0
	player.Userchairid = -1
	player.Usertableid = -1
	player.Usergamestatus = base.PlayerstatuWaitSitDown
	base.SendRspMsg(player, leaveres)
	if gametable.SitdownPlayers <= 0 {
		gametable.ResetTable()
	} else {
		player.Resetgamedata()
	}
	//记录离桌时间点
	player.Userlastopertime = time.Now()
	f.CallBackLeaveTable(player, gametable)
}
func (f *FactoryGameLogic) handleLoginOut(args []interface{}) {
	playerInt := args[1].(base.IPlayerNode)
	clientReq := args[0].(*msg.RequestData)
	var player *base.ClientNode = nil
	if playerInt.IsProxyNode() {
		proxyNode := playerInt.(*base.ProxyNode)
		proxyplayer, ok := proxyNode.MapClient[clientReq.ReqID]
		if !ok {
			return
		} else {
			player = proxyplayer
		}
	} else {
		player = playerInt.(*base.ClientNode)
	}
	f.OnPlayerClose(player)
}
func (f *FactoryGameLogic) handleloginout(player base.IPlayerNode) {
	playernode := player.(*base.ClientNode)
	logout := msg.Accountwdatareq{
		UserID:         playernode.Usernodeinfo.Userid,
		GameID:         conf.RoomInfo.GameID,
		GameCoin:       playernode.Useraccountdbw.Gamecoin,
		PrizeTicket:    playernode.Useraccountdbw.Goldbean,
		GameTax:        playernode.Useraccountdbw.Gametax,
		GamePlaytime:   playernode.Useraccountdbw.Gameplaytime,
		GameOnlinetime: int32(time.Now().Unix() - playernode.Userlogintime.Unix()),
		GameWintimes:   playernode.Useraccountdbw.Gamewintimes,
		GameLosetimes:  playernode.Useraccountdbw.Gamelosetimes,
		GameCoinwin:    playernode.Useraccountdbw.Gamecoinwin,
	}
	for _, propidlist := range playernode.Useraccountdbw.Proplist {
		for _, prop := range propidlist {
			logout.PropList = append(logout.PropList, prop)
		}
	}
	//只有成功验证后，才可以更新在线数据
	if playernode.Usergamestatus > base.PlayerstatuWaitAuthen {
		db.SaveUserProperty(logout) //更新积分等属性
		db.SaveGameData(logout)
	}
	db.DeleteOnline(playernode.Usernodeinfo.Userid)              //清除在线
	base.PlayerList.DeletePlayer(playernode.Usernodeinfo.Userid) //删除内存管理节点
}
