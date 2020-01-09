package robotlogic

import (
	"strings"
	"sync/atomic"

	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/msg"

	"github.com/name5566/leaf/log"
)

type RobotLogic struct {
	MapReqHandler map[string]base.RobotHandler
	freeRobotList []*RobotNode
	workingRobot  map[int64]*RobotNode
}
type RobotNode struct {
	*base.PlayerNode
	UserID   int64
	NickName string
	GameCoin int64
	TableID  int32
	ChairID  int32
}

var gRobotCalID uint64 = 1

func (p *RobotNode) Initialize() {
	p.PlayerID = atomic.AddUint64(&gRobotCalID, 1)
	p.TableID = -1
	p.ChairID = -1
}
func (p *RobotNode) GameBegin() {}
func (p *RobotNode) GameEnd()   {}
func (f *RobotLogic) HandleRobotMsg(args []interface{}) {
	msg := args[0].(*msg.RobotMessage)
	routeParam := strings.Split(msg.Route, ".")
	if fn, ok := f.MapReqHandler[routeParam[len(routeParam)-1]]; ok {
		fn(msg.UserID, msg.ReqMsg)
	} else {
		log.Error("route:%s not found", msg.Route)
	}
}
func (f *RobotLogic) OnCreateRobot() base.IPlayerNode {
	return &RobotNode{}
}
func (f *RobotLogic) RegisteRobotMsg() {
	f.MapReqHandler = make(map[string]base.RobotHandler)
	f.workingRobot = make(map[int64]*RobotNode)
	//此处添加机器人消息处理句柄
	f.MapReqHandler["LoginRes"] = f.handleApplyRobot
	f.MapReqHandler["SitDownRes"] = f.handleSitDownRes
	f.MapReqHandler["TablePlayerLeaveRes"] = f.handleLeaveRes
}
func (f *RobotLogic) handleApplyRobot(userID int64, robotmsg interface{}) {
	loginRes := robotmsg.(msg.ApplyRobotReq)
	if _, ok := f.workingRobot[loginRes.UserID]; ok {
		log.Error("robot:%v is already in use", loginRes.UserID)
		return
	}
	player := f.OnCreateRobot().(*RobotNode)
	player.Initialize()
	player.UserID = loginRes.UserID
	player.GameCoin = loginRes.GameCoin
	player.NickName = loginRes.NickName
	player.TableID = loginRes.TableID
	player.ChairID = loginRes.ChairID
	f.workingRobot[player.UserID] = player
	base.SendRobotMsg2Game("Game.SitDown", player.UserID, msg.SitDownReq{
		Tableid: player.TableID,
		Chairid: player.ChairID,
	})
}
func (f *RobotLogic) handleSitDownRes(userID int64, robotmsg interface{}) {
	sitRes := robotmsg.(msg.SitDownRes)
	if sitRes.Errcode != 0 {
		log.Error("robot:%v sitRes.Errcode:%v", userID, sitRes.Errcode)
	}
	if player, ok := f.workingRobot[userID]; ok {
		player.TableID = sitRes.Tableid
		player.ChairID = sitRes.Chairid
		base.SendRobotMsg2Game("Game.HandUp", player.UserID, msg.HandUpReq{})
	}
}
func (f *RobotLogic) handleLeaveRes(userID int64, robotmsg interface{}) {
	if player, ok := f.workingRobot[userID]; ok {
		log.Debug("robot:%v tableid:%v chairid:%v leavetable", userID, player.TableID, player.ChairID)
		delete(f.workingRobot, userID)
	}
}
func (f *RobotLogic) OnRobotLoginIn(player base.IPlayerNode, loginmsg interface{}) {

}
func (f *RobotLogic) OnRobotLoginOut(player base.IPlayerNode) {

}
func (f *RobotLogic) CloseRobot(player base.IPlayerNode) {
	playerNode := player.(*RobotNode)
	base.SendRobotMsg2Game("Game.LeaveScene", playerNode.UserID, msg.TablePlayerLeaveReq{})
}
func (f *RobotLogic) OnDestroy() {

}
