package base

type ITable interface {
	Init(chairnum int, flogic IGameLogic)
	ResetTable()
	GameBegin(gamenum int64)
	GameEnd()
}
