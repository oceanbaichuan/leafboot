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
	mapReqHandler           map[string]base.MsgHandler
	mapRobotHandler         map[string]base.MsgHandler
	mapFrontConnServer      map[string]*base.ProxyNode
	mapMiddlePlatConnServer map[string]*base.MiddlePlatNode
)

func init() {
	mapFrontConnServer = make(map[string]*base.ProxyNode)
	mapMiddlePlatConnServer = make(map[string]*base.MiddlePlatNode)
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
	f.AppMsgCallBackInit(mapReqHandler)
}
func (f *FactoryGameLogic) RegisteLoginMsgCallback(skeleton *module.Skeleton) {
	//skeleton.RegisterChanRPC(reflect.TypeOf(&msg.Logingamereq{}), f.handleLoginGameReq)
}
func (f *FactoryGameLogic) handleRequestMsg(args []interface{}) {
	req := args[0].(*msg.RequestData)
	// log.Debug("handleAdverprizereq req:%v", req)
	// // 消息的发送者
	playerinf := args[1].(gate.Agent).UserData().(base.IPlayerNode)
	routeParam := strings.Split(req.Route, ".")
	if fn, ok := mapReqHandler[routeParam[len(routeParam)-1]]; ok {
		args[1] = playerinf
		fn(args)
	} else {
		base.SendFailMsgWithID(playerinf, req.ReqID, 404, fmt.Sprintf("method:%s undefined", req.Route), nil)
		log.Error("method:%s undefined", req.Route)
	}
}

func (f *FactoryGameLogic) handleResponseMsg(args []interface{}) {

}
func (f *FactoryGameLogic) handleRobotMsg(args []interface{}) {
	req := args[0].(*msg.RobotMessage)
	routeParam := strings.Split(req.Route, ".")
	if player, ok := base.PlayerList.GetPlayer(req.UserID); ok {
		args = append(args, player)
	} else {
		log.Error("robotID:%v not created. method:%s", req.UserID, req.Route)
	}
	if fn, ok := mapReqHandler[routeParam[len(routeParam)-1]]; ok {
		fn(args)
	} else {
		log.Error("robotID:%v method:%s undefined", req.UserID, req.Route)
	}
}
