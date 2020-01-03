package gate

import (
	"github.com/hudgit2019/leafboot/module/game"
	"github.com/hudgit2019/leafboot/msg"
)

func init() {
	// 这里指定消息 Hello 路由到 game 模块
	// 模块间使用 ChanRPC 通讯，消息路由也不例外
	// msg.Processor.SetRouter(&msg.Logingamereq{}, login.ChanRPC)
	// msg.Processor.SetRouter(&msg.MqfromManager{}, game.ChanRPC)
	// msg.Processor.SetRouter(&msg.Sitdownreq{}, game.ChanRPC)
	// msg.Processor.SetRouter(&msg.Handupreq{}, game.ChanRPC)
	// msg.Processor.SetRouter(&msg.Tableplayerleavereq{}, game.ChanRPC)
	// 模块间使用 ChanRPC 通讯，消息路由也不例外
	msg.Processor.SetRouter(&msg.RequestData{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.ResponseData{}, game.ChanRPC)
	//appmain.SetRouter()
}
