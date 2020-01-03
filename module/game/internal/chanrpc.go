package internal

import (
	"unsafe"

	"github.com/hudgit2019/leafboot/base"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

type iface struct {
	itype  uintptr
	ivalue uintptr
}

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
}
func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	player := GameLogic.OnCreatePlayer(a.RemoteAddr().String())
	player.Initialize()
	if player.IsProxyNode() {
		player.(*base.ProxyNode).Netagent = a
	} else {
		player.(*base.ClientNode).Netagent = a
	}
	a.SetUserData(player)
	GameLogic.OnPlayerConnect(player)
	log.Debug("playernode:%v  ip:%v connected", (*iface)(unsafe.Pointer(&a)), a.RemoteAddr())
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.Debug("%v closed playernode: %v", a.RemoteAddr(), (*iface)(unsafe.Pointer(&a)))
	if a.UserData() != nil {
		playernode := a.UserData().(base.IPlayerNode)
		GameLogic.OnPlayerClose(playernode)
	}
}
