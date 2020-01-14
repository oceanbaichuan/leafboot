package base

type ITable interface {
	Init(tableid int32, chairnum int32, flogic IGameLogic)
	ResetTable()
	GameBegin(gamenum int64)
	GameEnd()
	SetCustomData(data interface{})
	CustomData() interface{}
}
