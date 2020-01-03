package msg

func init() {
	// 这里我们注册了一个 JSON 消息 Hello
	Processor.Register(&Register2manager{})
	Processor.Register(&MqfromManager{})
}

type Register2manager struct {
	Gameid int
}

type MqfromManager struct {
	Event       string
	Idevent     int
	Gameid      int
	Levelgame   int
	Roomid      int
	Areaid      int
	Channelid   int
	Distinc_id  int
	Userlcid    string
	Saved_field interface{}
	Properties  interface{}
}
