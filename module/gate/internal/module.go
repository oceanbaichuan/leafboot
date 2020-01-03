package internal

import (
	"github.com/hudgit2019/leafboot/conf"
	"github.com/hudgit2019/leafboot/module/game"
	"github.com/hudgit2019/leafboot/msg"

	"github.com/name5566/leaf/gate"
)

type Module struct {
	*gate.Gate
}

func (m *Module) OnInit() {
	m.Gate = &gate.Gate{
		MaxConnNum:      conf.Server.MaxConnNum,
		PendingWriteNum: conf.PendingWriteNum,
		MaxMsgLen:       conf.MaxMsgLen,
		WSAddr:          conf.Server.WSAddr,
		HTTPTimeout:     conf.HTTPTimeout,
		CertFile:        conf.Server.CertFile,
		KeyFile:         conf.Server.KeyFile,
		TCPAddr:         conf.Server.TCPAddr,
		LenMsgLen:       conf.LenMsgLen,
		LittleEndian:    conf.LittleEndian,
		Processor:       msg.Processor,
		AgentChanRPC:    game.ChanRPC,
	}
	conf.NewTCPAgent = m.Gate.NewTCPAgent
	if err := game.Start(m.Gate); err != nil {
		panic(err)
	}
}

func (m *Module) OnDestroy() {
}
