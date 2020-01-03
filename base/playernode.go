package base

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/hudgit2019/leafboot/conf"
	"github.com/hudgit2019/leafboot/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/timer"
)

const (
	Playertimer_checkhandup = iota //入座举手定时器
	Playertimer_checkplay          //游戏对局定时器
)

var GPlayerNodeID uint64 = 1

type PlayerNode struct {
	Netagent    gate.Agent           //socket链接点
	Usertimemap map[int]*timer.Timer //定时器
	PlayerID    uint64               //节点序列号
}
type ClientNode struct {
	PlayerNode
	Usernodeinfo          msg.LoginRes
	Useraccountdbw        Userplaygamedata //可更新入用户库数据
	Userlogindeviceid     string
	Userlogindevicetype   int
	Userloginchannellong  string
	Userloginchannelshort int
	Userloginfrom         string
	Userloginsiteid       int
	Usergamestatus        int       //游戏状态
	Useroffline           bool      //是否掉线
	Userisrobot           bool      //是否机器人
	Userlastopertime      time.Time //上次操作时间,只针对登录，入座，准备,打牌,离桌做记录
	Userlogintime         time.Time //登录时间点
	Usertableid           int       //座位号
	Userchairid           int       //椅子号
	ProxyClientID         uint64    //代理客户端ID
	LoginIP               string    //客户端IP
}

type ProxyNode struct {
	PlayerNode
	MapClient map[uint64]*ClientNode
}

//IsProxyNode 本身是否为代理节点
func (player *PlayerNode) IsProxyNode() bool {
	return false
}

//IsProxyedNode 本身是否被代理节点
func (player *PlayerNode) IsProxyedNode() bool {
	return false
}
func (player *PlayerNode) Initialize() {
	player.PlayerID = atomic.AddUint64(&GPlayerNodeID, 1)
	player.Netagent = nil
	player.Usertimemap = make(map[int]*timer.Timer)
}
func (player *PlayerNode) GameBegin() {
	for k, v := range player.Usertimemap {
		if v != nil {
			v.Stop()
		}
		player.Usertimemap[k] = nil
	}
}

func (player *PlayerNode) GameEnd() {
	for k, v := range player.Usertimemap {
		if v != nil {
			v.Stop()
		}
		player.Usertimemap[k] = nil
	}
	player.SetTimer(Playertimer_checkhandup, time.Duration(conf.RoomInfo.SitNoHandUpCheckTime)*time.Second)
}
func (player *PlayerNode) Resetgamedata() {
	for k, v := range player.Usertimemap {
		if v != nil {
			v.Stop()
		}
		player.Usertimemap[k] = nil
	}
}

func (player *PlayerNode) KillTimer(timerid int) {
	if t, ok := player.Usertimemap[timerid]; ok {
		if t != nil {
			t.Stop()
		}
		player.Usertimemap[timerid] = nil
	}
}
func (player *PlayerNode) OnTimer() {

}
func (player *PlayerNode) SetTimer(timerid int, d time.Duration) (*timer.Timer, error) {
	if t, ok := player.Usertimemap[timerid]; ok {
		return t, errors.New("timerid already in use")
	}
	t := Skeleton.AfterFunc(d, player.OnTimer)
	player.Usertimemap[timerid] = t
	return t, nil
}
func (player *PlayerNode) HandleAutoGame() {

}
func (player *ClientNode) Initialize() {
	player.PlayerNode.Initialize()
	player.Usertableid = -1
	player.Userchairid = -1
	player.ProxyClientID = 0
	player.Userisrobot = false
	player.Usergamestatus = PlayerstatuInitial
	player.Usernodeinfo = msg.LoginRes{}
	player.Useraccountdbw.Proplist = make(map[int]map[int]msg.UserPropChange)
}
func (player *ClientNode) OnTimer() {
	//sitdown too long to handup,close or autohandup
	if player.Usergamestatus == PlayerstatuHaveSitDown {
		kickout := msg.Kickoutres{
			Errcode: msg.Kickout_toolong_handup,
			Errmsg:  "Kickout_toolong_handup",
		}
		SendRspMsg(player, kickout)
		ClosePlayer(player)
		delete(player.Usertimemap, Playertimer_checkhandup)
	} else if player.Usergamestatus >= PlayerstatuBeginGame {
		player.HandleAutoGame()
		delete(player.Usertimemap, Playertimer_checkplay)
	}
}

func (player *ClientNode) SyncOldPlayerData(playerold IPlayerNode) {
	//first sync delta to usernodeinfo
	oldplayer := playerold.(*ClientNode)
	player.Usernodeinfo.GameCoin += oldplayer.Useraccountdbw.Gamecoin
	player.Usernodeinfo.GoldBean += oldplayer.Useraccountdbw.Goldbean
	player.Usernodeinfo.GameCoinLose += oldplayer.Useraccountdbw.Gamecoinlose
	player.Usernodeinfo.GameCoinWin += oldplayer.Useraccountdbw.Gamecoinwin
	player.Usernodeinfo.GameWinTimes += oldplayer.Useraccountdbw.Gamewintimes
	player.Usernodeinfo.GameLoseTimes += oldplayer.Useraccountdbw.Gamelosetimes
	player.Usernodeinfo.GameOnlineTime += oldplayer.Useraccountdbw.Gameonlinetime
	player.Usernodeinfo.GamePlayTime += oldplayer.Useraccountdbw.Gameplaytime
	for propid, typelist := range oldplayer.Useraccountdbw.Proplist {
		for proptype, prop := range typelist {
			var bhasprop = false
			for i, selfprop := range player.Usernodeinfo.Proplist {
				if selfprop.Propid == propid && selfprop.Proptype == proptype {
					bhasprop = true
					if prop.Propnum != 0 {
						player.Usernodeinfo.Proplist[i].Propnum += prop.Propnum
					} else if prop.Deadline != 0 {
						player.Usernodeinfo.Proplist[i].Deadline = time.Unix(player.Usernodeinfo.Proplist[i].Deadline.Unix()+prop.Deadline, 0)
					}
					break
				}
			}
			if !bhasprop {
				propinfo := msg.PropInfo{
					Propid:   prop.Propid,
					Proptype: prop.Proptype,
					Propnum:  prop.Propnum,
					Deadline: time.Unix(time.Now().Unix()+prop.Deadline, 0),
				}
				player.Usernodeinfo.Proplist = append(player.Usernodeinfo.Proplist, propinfo)
			}
		}
	}
	//second sync wdbdata to newplayer
	player.Useraccountdbw = oldplayer.Useraccountdbw
}

//IsProxyedNode 本身是否被代理节点
func (player *ClientNode) IsProxyedNode() bool {
	return player.ProxyClientID > 0
}
func (player *ProxyNode) Initialize() {
	player.PlayerNode.Initialize()
	player.MapClient = make(map[uint64]*ClientNode)
}

//IsProxyNode 本身是否为代理节点
func (player *ProxyNode) IsProxyNode() bool {
	return true
}
