package appmain

import (
	"github.com/hudgit2019/leafboot/appmain/gamelogic"
	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/msg"

	"github.com/name5566/leaf/gate"
)

/*主玩法逻辑示例*/
const (
	PlayerstatuThinking       = base.PlayerstatuBeginGame + 1
	PlayerstatuWaitOtherGuess = base.PlayerstatuBeginGame + 2
)

var (
	guesschoice = [3]int8{1, 2, 3}
)

type GuessLogic struct {
	gamelogic.FactoryGameLogic
}

func (g *GuessLogic) CreateClientPlayer() base.IPlayerNode {
	return &GuessPlayerNode{}
}

//AppMsgCallBack 应用层用来注册消息回调，重写方法,例如:mapMsg["LeaveTable"]= f.handleLeaveTableReq
func (g *GuessLogic) AppMsgCallBackInit(mapMsg *map[string]base.MsgHandler) {
	//mapMsg["LeaveTable"]= f.handleLeaveTableReq
}
func (g *GuessLogic) handleGuessReq(args []interface{}) {
	a := args[1].(gate.Agent)
	player := a.UserData().(*GuessPlayerNode)
	if player.Usergamestatus != PlayerstatuThinking {
		return
	}
	req := args[0].(*GuessReq)
	var bValidtype = false
	for _, v := range guesschoice {
		if v == req.Guesstype {
			bValidtype = true
			break
		}
	}
	if !bValidtype {
		kickout := msg.Kickoutres{
			Errcode: msg.Kickout_room_closed + 1,
			Errmsg:  "invalid guess",
		}
		base.SendRspMsg(player, kickout)
		base.ClosePlayer(player)
		return
	}
	player.Usergamestatus = PlayerstatuWaitOtherGuess
	player.Guesstype = req.Guesstype
	player.KillTimer(base.Playertimer_checkplay)
	gussres := GuessRes{
		Guesstype: req.Guesstype,
	}
	base.SendRspMsg(player, gussres)
	tableintf, _ := g.GetTable(player.Usertableid)
	table := tableintf.(*base.GameTable)
	var overguesscount int32 = 0
	for _, v := range table.TablePlayers {
		if v.(*base.ClientNode).Usergamestatus == PlayerstatuWaitOtherGuess {
			overguesscount++
		}
	}
	//game over
	if overguesscount == table.ReadyPlayers {
		gussresult := GuessResult{}
		for i, v := range table.TablePlayers {
			//g.SavePlayerGameCoin(v, 100, 113),just sample,you should fill your own data to struct!
			g.SavePlayerGameEnd(v, base.Userplaygamedata{}, 112)
			g.WriteTableRoundLog(&msg.Playgamelog{})
			gussresult.Guesstype[i] = v.(*GuessPlayerNode).Guesstype
			gussresult.Socres[i] = 1000
		}

		for _, v := range table.TablePlayers {
			base.SendRspMsg(v, gussresult)
		}
		g.SaveTableGameEnd(table)
	}
}

//CallBackSendRoomInfo 游戏层重写，同步游戏特有信息到c端
func (g *GuessLogic) CallBackSendRoomInfo(player base.IPlayerNode) {

}

//CallBackSitDown 游戏层重写玩家入座成功后的逻辑
func (g *GuessLogic) CallBackSitDown(player base.IPlayerNode, table base.ITable) {
}

//CallLoginSuccess 游戏层重写登录成功后的逻辑
func (g *GuessLogic) CallLoginSuccess(player base.IPlayerNode) {
}

//CallBackLeaveTable 游戏层重写离桌后的逻辑
func (g *GuessLogic) CallBackLeaveTable(player base.IPlayerNode, table base.ITable) {
}

//CallBackOffline 游戏层重写用户断线的逻辑
func (g *GuessLogic) CallBackOffline(player base.IPlayerNode) {
}

//CallBackLogOut 游戏层重写用户完全退出房间后的逻辑
func (g *GuessLogic) CallBackLogOut(player base.IPlayerNode) {

}

//CallBackHandUp 游戏层重写玩家举手后的逻辑
func (g *GuessLogic) CallBackHandUp(player base.IPlayerNode, table base.ITable) {
}

//CallBackGameStart 游戏层重写游戏对局开始的逻辑
func (g *GuessLogic) CallBackGameStart(table base.ITable) {
	tableitem := table.(*base.GameTable)
	for _, playerintf := range tableitem.TablePlayers {
		player := playerintf.(*base.ClientNode)
		player.Usergamestatus = PlayerstatuThinking
		guessstart := GuessStartNotice{}
		base.SendRspMsg(player, guessstart)
		player.SetTimer(base.Playertimer_checkplay, 10)
	}
}

//CallBackLoginAgain 游戏层重写玩家断线重连的逻辑
func (g *GuessLogic) CallBackLoginAgain(player base.IPlayerNode) {
}

//
func (g *GuessLogic) AutoPlay(player base.IPlayerNode) {

}
func (g *GuessLogic) OnExit() {

}
