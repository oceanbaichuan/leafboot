package gameboot

import (
	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/conf"
	"github.com/hudgit2019/leafboot/module/game"
	"github.com/hudgit2019/leafboot/module/gate"
	"github.com/hudgit2019/leafboot/module/robot"

	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
)

func StartGame(fLogic base.IGameLogic, fRobot base.IRobot) error {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.BuringPointLogPath = conf.Server.BuringPointLogPath
	lconf.LogSplitHour = int(conf.Server.LogSplitHour)
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = int(conf.Server.ConsolePort)
	lconf.ProfilePath = conf.Server.ProfilePath
	//绑定游戏逻辑层
	base.GameChanRPC = game.ChanRPC
	base.RobotChanRPC = robot.ChanRPC
	game.InitLogic(fLogic)
	robot.InitLogic(fRobot)
	leaf.Run(
		robot.Module,
		game.Module,
		gate.Module,
	)
	return nil
}
