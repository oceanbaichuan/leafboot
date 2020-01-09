package base

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"time"

	"github.com/hudgit2019/leafboot/msg"

	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/module"
)

var (
	Skeleton   *module.Skeleton
	PlayerList *PlayerNodeList
	//Curgamenum   int64 //当前牌局计数，初始化为时间戳
	//GameTables   []ITable
	NetGate *gate.Gate
	//GameLogic    IGameLogic
	GameChanRPC  *chanrpc.Server
	RobotChanRPC *chanrpc.Server
)

type MsgHandler func(agrs []interface{})
type RobotHandler func(userID int64, msg interface{})

//CheckDateDiff srcTime:比较时间 destime:被比较时间  compType：day,month,year cmpValue:差值最大
func CheckDateDiff(srcTime time.Time, destime time.Time, compType string, cmpValue int) bool {
	switch compType {
	case "day":
		{
			if srcTime.YearDay()+cmpValue <= destime.YearDay() {
				return true
			}
		}
	case "month":
		{
			if int(srcTime.Month())+cmpValue <= int(destime.Month()) {
				return true
			}
		}
	case "year":
		{
			if srcTime.Year()+cmpValue <= destime.Year() {
				return true
			}
		}
	}
	return false
}
func init() {
	PlayerList = new(PlayerNodeList)
	PlayerList.Init()
}
func GernateToken(account string, userID int64) string {
	hashdig := md5.New()
	hashdig.Write([]byte(fmt.Sprintf("%s:%v:%v",
		account, userID, time.Now().Unix())))
	sercode := fmt.Sprintf("%x", hashdig.Sum(([]byte(""))))
	return sercode
}

//GameHandler 消息回调注入信道
func GameHandler(m interface{}, h interface{}) {
	Skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

// func GetTable(tableid int) (ITable, error) {
// 	if tableid < 1 || tableid > conf.RoomInfo.MaxTableNum {
// 		return nil, errors.New("table index out of range")
// 	}
// 	return GameTables[tableid-1], nil
// }
func ClosePlayer(player IPlayerNode) {
	playernode := player.(*PlayerNode)
	if playernode.Netagent != nil {
		playernode.Netagent.Close()
	}
}

// func Sendmsg(player IPlayerNode, args ...interface{}) error {
// 	if msg.Processor.Type() == json.Processor_Type {
// 		netagent := player.(*PlayerNode).Netagent
// 		rspBody := &msg.Responsedata{}
// 		rspBody.Rspname = reflect.TypeOf(args[0]).Name()
// 		log.Debug("rspBody.Rspname %v", rspBody.Rspname)
// 		rspBody.Rspdata = args[0]
// 		netagent.WriteMsg(rspBody)
// 	} else if msg.Processor.Type() == protobuf.Processor_Type {
// 		netagent := player.(*PlayerNode).Netagent
// 		if len(args) < 2 {
// 			log.Error("arguments must have msg, msgtype")
// 			return errors.New("arguments must have msg, msgtype")
// 		}
// 		netagent.Write8HeaderMsg(args[1].(uint16), args[0])
// 	} else {
// 		log.Error("Sendmsg msg.Processor is invalid:%v", reflect.TypeOf(msg.Processor))
// 		return errors.New("Sendmsg msg.Processor is invalid")
// 	}
// 	return nil
// }

type IGameLogic interface {
	//可重写频率低接口
	Start(netgate *gate.Gate) error
	Gamestart(table ITable)
	OnCreateTable() ITable
	OnCreatePlayer(addr string) IPlayerNode
	RegisteGameMsgCallback(skeleton *module.Skeleton)
	RegisteLoginMsgCallback(skeleton *module.Skeleton)
	OnPlayerConnect(player IPlayerNode)
	OnPlayerClose(player IPlayerNode)
	ClosePlayer(player IPlayerNode)
	SavePlayerGameCoin(player IPlayerNode, gamecoin int64, writesource int32)
	SavePlayerGoldBean(player IPlayerNode, goldbean int32, writesource int32)
	SavePlayerProp(player IPlayerNode, propinfo msg.UserPropChange, writesource int32)
	SavePlayerGameEnd(player IPlayerNode, datachanged Userplaygamedata)
	SaveTableGameEnd(table ITable)
	WriteLoginRoomLog(loginlog interface{})
	WriteLeaveRoomLog(leavelog interface{})
	WriteTableRoundLog(playlog interface{})
	WriteAttributionLog(attrlog interface{})

	//可大量重写接口
	CallBackSendRoomInfo(player IPlayerNode)
	CallBackSitDown(player IPlayerNode, table ITable)
	CallLoginSuccess(player IPlayerNode)
	CallBackLeaveTable(player IPlayerNode, table ITable)
	CallBackOffline(player IPlayerNode)
	CallBackLogOut(player IPlayerNode)
	CallBackHandUp(player IPlayerNode, table ITable)
	CallBackGameStart(table ITable)
	CallBackLoginAgain(player IPlayerNode)
	AutoPlay(player IPlayerNode)
	OnDestroy()
}

type IRobot interface {
	HandleRobotMsg(args []interface{})
	OnCreateRobot() IPlayerNode
	RegisteRobotMsg()
	OnRobotLoginIn(player IPlayerNode, loginmsg interface{})
	OnRobotLoginOut(player IPlayerNode)
	CloseRobot(player IPlayerNode)
	OnDestroy()
}
