package msg

type RobotMessage struct {
	Route  string //路由
	UserID int64  //用户ID
	Msg    interface{}
}
