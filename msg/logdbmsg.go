package msg

import (
	"time"

	"go.uber.org/zap/zapcore"
)

//Logingamelog 进入房间日志
type Logingamelog struct {
	time       int64  //事件发生时间(毫秒)
	user_id    int64  //玩家id
	device_id  string //设备id
	project    string //项目名：datatalk
	event_type string //事件类型：核心事件
	event      string //事件名称：进入房间

	model      string //设备型号
	os         string //操作系统
	os_version string //系统版本
	network    string //网络类型
	ip         string //ip地址

	lobby_version     string //大厅版本
	site_id           int    //站点id
	market_channel_id int    //市场渠道id
	app_channel_info  string //长渠道id
	channel_id        int    //短渠道id

	game_id int //游戏id
	room_id int //房间id

}

func (logingamelog *Logingamelog) MarshalLogObject(objencoder zapcore.ObjectEncoder) error {
	objencoder.AddInt64("time", time.Now().UnixNano()/1e6)
	objencoder.AddInt64("user_id", logingamelog.user_id)
	objencoder.AddString("device_id", logingamelog.device_id)
	objencoder.AddString("project", "datatalk")
	objencoder.AddString("event_type", "kernel")
	objencoder.AddString("event", "room_enter")
	objencoder.AddString("model", logingamelog.model)
	objencoder.AddString("os", logingamelog.os)
	objencoder.AddString("os_version", logingamelog.os_version)
	objencoder.AddString("network", logingamelog.network)
	objencoder.AddString("ip", logingamelog.ip)
	objencoder.AddString("lobby_version", logingamelog.lobby_version)
	objencoder.AddInt("site_id", logingamelog.site_id)
	objencoder.AddInt("market_channel_id", logingamelog.market_channel_id)
	objencoder.AddString("app_channel_info", logingamelog.app_channel_info)
	objencoder.AddInt("channel_id", logingamelog.channel_id)
	objencoder.AddInt("game_id", logingamelog.game_id)
	objencoder.AddInt("room_id", logingamelog.room_id)
	return nil
}

//Leavegamelog 退出房间日志
type Leavegamelog struct {
	time              int    //事件发生时间(毫秒)
	user_id           int    //玩家id
	device_id         string //设备id
	project           string //项目名：datatalk
	event_type        string //事件类型：核心事件
	event             string //事件名称：离开房间
	model             string //设备型号
	os                string //操作系统
	os_version        string //系统版本
	network           string //网络类型
	ip                string //ip地址
	lobby_version     string //大厅版本
	site_id           int    //站点id
	market_channel_id int    //市场渠道id
	app_channel_info  string //长渠道id
	channel_id        int    //短渠道id
	game_id           int    //游戏id
	room_id           int    //房间id

}

func (leavegamelog *Leavegamelog) MarshalLogObject(objencoder zapcore.ObjectEncoder) error {
	objencoder.AddInt64("time", time.Now().UnixNano()/1e6)
	objencoder.AddInt("user_id", leavegamelog.user_id)
	objencoder.AddString("device_id", leavegamelog.device_id)
	objencoder.AddString("project", "datatalk")
	objencoder.AddString("event_type", "kernel")
	objencoder.AddString("event", "room_leave")
	objencoder.AddString("model", leavegamelog.model)
	objencoder.AddString("os", leavegamelog.os)
	objencoder.AddString("os_version", leavegamelog.os_version)
	objencoder.AddString("network", leavegamelog.network)
	objencoder.AddString("ip", leavegamelog.ip)
	objencoder.AddString("lobby_version", leavegamelog.lobby_version)
	objencoder.AddInt("site_id", leavegamelog.site_id)
	objencoder.AddInt("market_channel_id", leavegamelog.market_channel_id)
	objencoder.AddString("app_channel_info", leavegamelog.app_channel_info)
	objencoder.AddInt("channel_id", leavegamelog.channel_id)
	objencoder.AddInt("game_id", leavegamelog.game_id)
	objencoder.AddInt("room_id", leavegamelog.room_id)
	return nil
}

//Playgamelog 游戏对局日志
type Playgamelog struct {
	time       int    //事件发生时间(毫秒)
	user_id    int    //玩家id
	device_id  string //设备id
	project    string //项目名：datatalk
	event_type string //事件类型：核心事件
	event      string //事件名称：AB对局

	model      string //设备型号
	os         string //操作系统
	os_version string //系统版本
	network    string //网络类型
	ip         string //ip地址

	lobby_version     string //大厅版本
	site_id           int    //站点id
	market_channel_id int    //市场渠道id
	app_channel_info  string //长渠道id
	channel_id        int    //短渠道id

	game_id  int //游戏id
	room_id  int //房间id
	table_id int //桌子id

	launch_from   string //对局来源
	combat_serial string //对局编号
	time_begin    int    //对局开始时间(毫秒)
	time_end      int    //对局结束时间(毫秒)

	multiple        int //翻倍
	gamecoin_change int //积分变化
	gamecoin_remain int //积分剩余
	tax_num         int //税收

	is_freshhelp int //是否新手补助阶段
	is_banker    int //是否庄家
	is_bankrupt  int //是否破产
	is_roomup    int //是否升场
	is_roomdown  int //是否降场
	is_robot     int //是否机器人
	is_offline   int //是否离线

	cards_begin string //开始手牌
	cards_end   string //结束手牌

}

func (playgamelog *Playgamelog) MarshalLogObject(objencoder zapcore.ObjectEncoder) error {
	objencoder.AddInt64("time", time.Now().UnixNano()/1e6)
	objencoder.AddInt("user_id", playgamelog.user_id)
	objencoder.AddString("device_id", playgamelog.device_id)
	objencoder.AddString("project", "datatalk")
	objencoder.AddString("event_type", "kernel")
	objencoder.AddString("event", "room_leave")
	objencoder.AddString("model", playgamelog.model)
	objencoder.AddString("os", playgamelog.os)
	objencoder.AddString("os_version", playgamelog.os_version)
	objencoder.AddString("network", playgamelog.network)
	objencoder.AddString("ip", playgamelog.ip)
	objencoder.AddString("lobby_version", playgamelog.lobby_version)
	objencoder.AddInt("site_id", playgamelog.site_id)
	objencoder.AddInt("market_channel_id", playgamelog.market_channel_id)
	objencoder.AddString("app_channel_info", playgamelog.app_channel_info)
	objencoder.AddInt("channel_id", playgamelog.channel_id)
	objencoder.AddInt("game_id", playgamelog.game_id)
	objencoder.AddInt("room_id", playgamelog.room_id)
	objencoder.AddInt("table_id", playgamelog.table_id)
	objencoder.AddString("launch_from", playgamelog.launch_from)
	objencoder.AddString("combat_serial", playgamelog.combat_serial)
	objencoder.AddInt("time_begin", playgamelog.time_begin)
	objencoder.AddInt("time_end", playgamelog.time_end)
	objencoder.AddInt("multiple", playgamelog.multiple)
	objencoder.AddInt("gamecoin_change", playgamelog.gamecoin_change)
	objencoder.AddInt("gamecoin_remain", playgamelog.gamecoin_remain)
	objencoder.AddInt("tax_num", playgamelog.tax_num)
	objencoder.AddInt("is_freshhelp", playgamelog.is_freshhelp)
	objencoder.AddInt("is_banker", playgamelog.is_banker)
	objencoder.AddInt("is_bankrupt", playgamelog.is_bankrupt)
	objencoder.AddInt("is_roomup", playgamelog.is_roomup)
	objencoder.AddInt("is_roomdown", playgamelog.is_roomdown)
	objencoder.AddInt("is_robot", playgamelog.is_robot)
	objencoder.AddInt("is_offline", playgamelog.is_offline)
	objencoder.AddString("cards_begin", playgamelog.cards_begin)
	objencoder.AddString("cards_end", playgamelog.cards_end)
	return nil
}

//AttributeChangelog 用户货币道具等属性改变日志
type AttributeChangelog struct {
	time       int    //事件发生时间(毫秒)
	user_id    int    //玩家id
	device_id  string //设备id
	project    string //项目名：datatalk
	event_type string //事件类型：附加事件
	event      string //事件名称：属性变化

	model      string //设备型号
	os         string //操作系统
	os_version string //系统版本
	network    string //网络类型
	ip         string //ip地址

	lobby_version     string //大厅版本
	site_id           int    //站点id
	market_channel_id int    //市场渠道id
	app_channel_info  string //长渠道id
	channel_id        int    //短渠道id

	game_id       int    //游戏id 0:大厅
	room_id       int    //房间id
	table_id      int    //桌子id
	combat_serial string //对局编号(对局)
	order_serial  string //订单编号(充值/兑换)

	source_id        int //来源id
	attribute_type   int //属性类型
	attribute_id     int //属性id
	attribute_change int //属性变化
	attribute_remain int //属性剩余

	project_id  int //项目id
	activity_id int //活动id

}

func (proplog *AttributeChangelog) MarshalLogObject(objencoder zapcore.ObjectEncoder) error {
	objencoder.AddInt64("time", time.Now().UnixNano()/1e6)
	objencoder.AddInt("user_id", proplog.user_id)
	objencoder.AddString("device_id", proplog.device_id)
	objencoder.AddString("project", "datatalk")
	objencoder.AddString("event_type", "kernel")
	objencoder.AddString("event", "room_leave")
	objencoder.AddString("model", proplog.model)
	objencoder.AddString("os", proplog.os)
	objencoder.AddString("os_version", proplog.os_version)
	objencoder.AddString("network", proplog.network)
	objencoder.AddString("ip", proplog.ip)
	objencoder.AddString("lobby_version", proplog.lobby_version)
	objencoder.AddInt("site_id", proplog.site_id)
	objencoder.AddInt("market_channel_id", proplog.market_channel_id)
	objencoder.AddString("app_channel_info", proplog.app_channel_info)
	objencoder.AddInt("channel_id", proplog.channel_id)
	objencoder.AddInt("game_id", proplog.game_id)
	objencoder.AddInt("room_id", proplog.room_id)
	objencoder.AddInt("table_id", proplog.table_id)
	objencoder.AddString("order_serial", proplog.order_serial)
	objencoder.AddString("combat_serial", proplog.combat_serial)
	objencoder.AddInt("source_id", proplog.source_id)
	objencoder.AddInt("attribute_type", proplog.attribute_type)
	objencoder.AddInt("attribute_id", proplog.attribute_id)
	objencoder.AddInt("attribute_change", proplog.attribute_change)
	objencoder.AddInt("gamecoin_remain", proplog.attribute_remain)
	objencoder.AddInt("project_id", proplog.project_id)
	return nil
}
