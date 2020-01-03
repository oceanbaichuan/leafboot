package base

import (
	"errors"
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
		netagent.WriteMsg(rspBody)
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
		netagent.WriteMsg(m)
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
		netagent.WriteMsg(rspBody)
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
		netagent.WriteMsg(rspBody)
	}
}

//SendMsg2Robot 玩法区与AI区通信
func SendGameMsg2Robot(args []interface{}) error {
	if len(args) != 2 {
		return errors.New("Invalid arguments!first must be PlayerNode,second must be msg!")
	}
	return RobotChanRPC.Call0(reflect.TypeOf(args[1]), args...)
}

//SendMsg2GameLogic AI区与玩法区通信
func SendRobotMsg2Game(args []interface{}) error {
	if len(args) != 2 {
		return errors.New("Invalid arguments!first must be PlayerNode,second must be msg!")
	}
	return GameChanRPC.Call0(reflect.TypeOf(args[1]), args...)
}
