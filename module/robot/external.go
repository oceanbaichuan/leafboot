package robot

import (
	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/module/robot/internal"
)

var (
	Module  = new(internal.Module)
	ChanRPC = internal.ChanRPC
)

func InitLogic(flogic base.IRobot) {
	internal.RobotLogic = flogic
	if flogic != nil {
		internal.RobotLogic.InitAppMain(flogic)
	}
}
