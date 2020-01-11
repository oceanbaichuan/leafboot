package msg

//MqFromConSumer MQ消费数据
type MqFromConSumer struct {
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
