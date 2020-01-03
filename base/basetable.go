package base

type ITable interface {
	Init(chairnum int32, flogic IGameLogic)
	ResetTable()
	GameBegin(gamenum int64)
	GameEnd()
}
