package base

import (
	"reflect"

	"github.com/hudgit2019/leafboot/msg"

	"github.com/name5566/leaf/log"
)

func SendRspMsgWithID(player *ProxyNode, rspID uint64, m interface{}) {
	rspBody := &msg.ResponseData{}
	rspBody.RspName = reflect.TypeOf(m).Name()
	log.Debug("rspBody.Rspname %v", rspBody.RspName)
	rspBody.RspData = m
	rspBody.RspID = rspID
	rspBody.Code = 200
	rspBody.Message = "success"
	player.Netagent.WriteMsg(rspBody)
}
func SendRspMsg(player IPlayerNode, m interface{}) error {
	rspBody := &msg.ResponseData{}
	rspBody.RspName = reflect.TypeOf(m).Name()
	log.Debug("rspBody.Rspname %v", rspBody.RspName)
	rspBody.RspData = m
	rspBody.Code = 200
	rspBody.Message = "success"
	if player.IsProxyNode() {
		netagent := player.(*ProxyNode).Netagent
		netagent.WriteMsg(rspBody)
	} else {
		clientNode := player.(*ClientNode)
		netagent := clientNode.Netagent
		if clientNode.IsProxyedNode() {
			rspBody.RspID = clientNode.ProxyClientID
		} else {
			rspBody.RspID = clientNode.PlayerID
		}
		if !clientNode.Userisrobot {
			netagent.WriteMsg(rspBody)
		} else {
			rspBody.RspID = uint64(clientNode.Usernodeinfo.Userid)
			sendGameMsg2Robot(rspBody)
		}
	}
	return nil
}

//SendAutoReqMsg 自动包装消息内容
func SendAutoReqMsg(player IPlayerNode, route string, reqID uint64, m interface{}) error {
	reqBody := &msg.RequestData{
		Route:   route,
		ReqID:   reqID,
		ReqData: m,
	}
	if player.IsProxyNode() {
		netagent := player.(*ProxyNode).Netagent
		netagent.WriteMsg(reqBody)
	} else {
		clientNode := player.(*ClientNode)
		netagent := clientNode.Netagent
		if !clientNode.Userisrobot {
			netagent.WriteMsg(reqBody)
		} else {

		}
	}
	return nil
}
func SendReqMsg(player IPlayerNode, m interface{}) error {
	if player.IsProxyNode() {
		netagent := player.(*ProxyNode).Netagent
		netagent.WriteMsg(m)
	} else {
		clientNode := player.(*ClientNode)
		netagent := clientNode.Netagent
		if !clientNode.Userisrobot {
			netagent.WriteMsg(m)
		} else {
		}
	}
	return nil
}
func SendFailMsg(player IPlayerNode, code int32, message string, m interface{}) {
	rspBody := &msg.ResponseData{}
	rspBody.Code = code
	rspBody.Message = message
	if player.IsProxyNode() {
		netagent := player.(*ProxyNode).Netagent
		netagent.WriteMsg(rspBody)
	} else {
		clientNode := player.(*ClientNode)
		netagent := clientNode.Netagent
		if clientNode.IsProxyedNode() {
			rspBody.RspID = clientNode.ProxyClientID
		} else {
			rspBody.RspID = clientNode.PlayerID
		}
		if !clientNode.Userisrobot {
			netagent.WriteMsg(rspBody)
		} else {
			rspBody.RspID = uint64(clientNode.Usernodeinfo.Userid)
			sendGameMsg2Robot(rspBody)
		}

	}
}
func SendFailMsgWithID(player IPlayerNode, rspID uint64, code int32, message string, m interface{}) {
	rspBody := &msg.ResponseData{}
	rspBody.Code = code
	rspBody.Message = message
	rspBody.RspID = rspID
	if player.IsProxyNode() {
		netagent := player.(*ProxyNode).Netagent
		netagent.WriteMsg(rspBody)
	} else {
		netagent := player.(*ClientNode).Netagent
		if !player.(*ClientNode).Userisrobot {
			netagent.WriteMsg(rspBody)
		} else {
			rspBody.RspID = uint64(player.(*ClientNode).Usernodeinfo.Userid)
			sendGameMsg2Robot(rspBody)
		}
	}
}

//SendMsg2Robot 玩法区与AI区通信
func sendGameMsg2Robot(robotmsg interface{}) error {
	// if len(args) != 2 {
	// 	return errors.New("Invalid arguments!first must be PlayerNode,second must be msg!")
	// }
	return RobotChanRPC.Call0(reflect.TypeOf(&msg.ResponseData{}), robotmsg)
}

//SendMsg2GameLogic AI区与玩法区通信
func SendRobotMsg2Game(route string, userID int64, robotmsg interface{}) error {
	reqBody := &msg.RobotMessage{
		Route:  route,
		UserID: userID,
		ReqMsg: robotmsg,
	}
	return GameChanRPC.Call0(reflect.TypeOf(&msg.RobotMessage{}), reqBody)
}
