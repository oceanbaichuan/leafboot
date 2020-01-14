package game

import (
	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/module/game/internal"

	"github.com/name5566/leaf/gate"
)

var (
	Module  = new(internal.Module)
	ChanRPC = internal.ChanRPC
)

func InitLogic(f base.IGameLogic) {
	internal.GameLogic = f
	internal.GameLogic.InitAppMain(f)
}
func Start(gate *gate.Gate) error {
	return internal.GameLogic.Start(gate)
}
