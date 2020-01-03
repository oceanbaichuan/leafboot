package internal

import (
	"github.com/hudgit2019/leafboot/base"

	"github.com/name5566/leaf/module"
)

var (
	skeleton  = base.NewSkeleton()
	ChanRPC   = skeleton.ChanRPCServer
	GameLogic base.IGameLogic
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	GameLogic.RegisteGameMsgCallback(skeleton)
}
func (m *Module) OnDestroy() {
	GameLogic.OnDestroy()
}
