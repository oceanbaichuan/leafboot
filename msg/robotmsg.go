package msg

type RobotMessage struct {
	Route  string
	UserID int64
	ReqMsg interface{}
}
type ApplyRobotReq struct {
	UserID   int64
	NickName string
	GameCoin int64
	TableID  int32
	ChairID  int32
}
