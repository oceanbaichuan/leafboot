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
}

var gRobotCalID uint64 = 1

func (p *RobotNode) Initialize() {
	p.PlayerID = atomic.AddUint64(&gRobotCalID, 1)
}
func (p *RobotNode) GameBegin() {}
func (p *RobotNode) GameEnd()   {}
func (f *RobotLogic) HandleRobotMsg(args []interface{}) {
	userID := args[0].(int64)
	msg := args[1].(*msg.RobotMessage)
	routeParam := strings.Split(msg.Route, ".")
	if fn, ok := f.MapReqHandler[routeParam[len(routeParam)-1]]; ok {
		fn(userID, msg.Msg)
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
	f.MapReqHandler["ApplyRobot"] = f.handleApplyRobot
	f.MapReqHandler["SitDownRes"] = f.handleSitDownRes
	f.MapReqHandler["LeaveReq"] = f.handleLeaveReq
}
func (f *RobotLogic) handleApplyRobot(userID int64, msg interface{}) {

}
func (f *RobotLogic) handleSitDownRes(userID int64, msg interface{}) {

}
func (f *RobotLogic) handleLeaveReq(userID int64, msg interface{}) {

}
func (f *RobotLogic) OnRobotLoginIn(player base.IPlayerNode, loginmsg interface{}) {

}
func (f *RobotLogic) OnRobotLoginOut(player base.IPlayerNode) {

}
func (f *RobotLogic) CloseRobot(player base.IPlayerNode) {

}
func (f *RobotLogic) OnDestroy() {

}
