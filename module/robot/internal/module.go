package internal

import (
	"reflect"

	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/msg"

	"github.com/name5566/leaf/module"
)

var (
	skeleton   = base.NewSkeleton()
	ChanRPC    = skeleton.ChanRPCServer
	RobotLogic base.IRobot
)

type Module struct {
	*module.Skeleton
}

//GameHandler 消息回调注入信道
func gameHandler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}
func (m *Module) OnInit() {
	m.Skeleton = skeleton
	if RobotLogic != nil {
		gameHandler(&msg.RobotMessage{}, RobotLogic.HandleRobotMsg)
		RobotLogic.RegisteRobotMsg()
	}
}
func (m *Module) OnDestroy() {
	if RobotLogic != nil {
		RobotLogic.OnDestroy()
	}
}
