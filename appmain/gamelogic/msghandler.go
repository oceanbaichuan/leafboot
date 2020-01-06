package gamelogic

import (

	//"encoding/json"

	"fmt"
	"strings"

	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/module"
)

var (
	mapReqHandler   map[string]base.MsgHandler
	mapRobotHandler map[string]base.MsgHandler
	mapConnServer   map[string]*base.ProxyNode
)

func init() {
	mapConnServer = make(map[string]*base.ProxyNode)
	mapReqHandler = make(map[string]base.MsgHandler)
	mapRobotHandler = make(map[string]base.MsgHandler)

}

func (f *FactoryGameLogic) RegisteGameMsgCallback(skeleton *module.Skeleton) {
	base.Skeleton = skeleton
	//监听请求消息
	base.GameHandler(&msg.RequestData{}, f.handleRequestMsg)
	//监听反馈消息
	base.GameHandler(&msg.ResponseData{}, f.handleResponseMsg)
	base.GameHandler(&msg.RobotMessage{}, f.handleRobotMsg)
	//此处添加游戏框架消息句柄
	mapReqHandler["SitDown"] = f.handleSitdownReq
	mapReqHandler["HandUp"] = f.handleHandUpReq
	mapReqHandler["EnterScene"] = f.handleLoginGameReq
	mapReqHandler["LeaveScene"] = f.handleLoginOut
	mapReqHandler["LeaveTable"] = f.handleLeaveTableReq
	//此处添加机器人框架消息句柄
	f.AppMsgCallBackInit(&mapReqHandler)
}
func (f *FactoryGameLogic) RegisteLoginMsgCallback(skeleton *module.Skeleton) {
	//skeleton.RegisterChanRPC(reflect.TypeOf(&msg.Logingamereq{}), f.handleLoginGameReq)
}
func (f *FactoryGameLogic) handleRequestMsg(args []interface{}) {
	req := args[0].(*msg.RequestData)
	// log.Debug("handleAdverprizereq req:%v", req)
	// // 消息的发送者
	a := args[1].(gate.Agent)
	routeParam := strings.Split(req.Route, ".")
	if fn, ok := mapReqHandler[routeParam[len(routeParam)-1]]; ok {
		fn(args)
	} else {
		playerinf := a.UserData().(base.IPlayerNode)
		base.SendFailMsgWithID(playerinf, req.ReqID, 404, fmt.Sprintf("method:%s undefined", req.Route), nil)
		log.Error("method:%s undefined", req.Route)
	}
}

func (f *FactoryGameLogic) handleResponseMsg(args []interface{}) {

}
func (f *FactoryGameLogic) handleRobotMsg(args []interface{}) {

}
