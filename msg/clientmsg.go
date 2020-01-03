package msg

import "time"

//EnterScenceReq 进入业务服
type EnterScenceReq struct {
	Account string //账号
	Passwd  string //密码
	UserID  int64  //用户ID
	AppID   string //产品ID
	MacID   string //机器码
	LoginIP string //登录IP
}

//LeaveScenceReq 进入业务服
type LeaveScenceReq struct {
	UserID int64 //用户ID
}
type LeaveScenceRes struct {
	Route string //用户ID
}

//AskConnReq 反馈合适连接服务器
type EnterScenceRes struct {
	ConnAddr string
}

//游戏登录
const (
	LoginErr_tokenerr = iota
	LoginErr_nouser
	LoginErr_inotherroom
	LoginErr_gamecoin_notenough
	LoginErr_gamecoin_toomuch
	LoginErr_accountforbidden
)

// type Logingamereq struct {
// 	Userid            int
// 	Loginlongchannel  string
// 	Loginshortchannel int
// 	Logindeviceid     string
// 	Logindevicetype   int    //1:PC 2:安卓 3:IOS
// 	Loginsiteid       int    //登陆站点
// 	Loginfrom         string //登陆来源
// 	Logintoken        string //登陆token
// }
type LoginRes struct {
	Userid         int64
	Errcode        int
	Account        string    //用户账号
	AcIndex        int       //用户账号所属机器码的序号
	AcType         int       //账号类型
	Freezed        bool      //账号是否冻结
	NickName       string    //用户昵称
	GameCoin       int64     //积分
	BankCoin       int64     //
	GoldBean       int       //金豆
	AllGoldBean    int       //
	RegChan        int       //注册短渠道
	RegTime        time.Time //注册时间
	RegSiteID      int
	Gender         int8   //性别
	HeadID         int    //系统头像id
	ThirdHeadUrl   string //第三方头像地址
	Phonebinded    string //绑定手机号
	VipExp         int    //
	GameWinTimes   int    //
	GameLoseTimes  int    //
	GamePlayTime   int    //
	PlatPlayTime   int    //
	GameOnlineTime int    //
	PlatOnlineTime int    //
	GameCoinPlay   int64  //
	GameCoinWin    int64  //
	GameCoinLose   int64  //
	OtherGameID    int    //
	OtherRoomID    string //
	GameStatus     int
	Proplist       []PropInfo //道具列表
}

type TablePlayer struct {
	Userid        int64 //
	Siteid        int
	Nickname      string //
	Gamecoin      int64
	Bankcoin      int64
	Goldbean      int
	Sysheadid     int
	Thirdheadurl  string //
	Gender        int8
	Vipexp        int //
	Gamewintimes  int //
	Gamelosetimes int //
	Tableid       int
	Chairid       int //
	Gamestatus    int //
}

type Loginagainres struct {
	Selfbaseinfo LoginRes
	Tableplayers []TablePlayer
}

//入座
const (
	SitErr_NoFitLocation = iota
	SitErr_HaveSit
	SitErr_ForbidWatch
	SitErr_NoLogin
)

type Sitdownreq struct {
	Tableid int //
	Chairid int //
}

type Sitdownres struct {
	Errcode int //错误码
	Tableid int //桌子号
	Chairid int //椅子号
	Players []TablePlayer
}

type TableJoinPlayer struct {
	Player TablePlayer
}

//Handuperr
const (
	Handuperr_no_sitdown = iota
	Handuperr_already_hand
	Handuperr_coin_notenough
	Handuperr_coin_toomuch
	Handuperr_room_closed
)

type Handupreq struct {
	Handup int8
}

type Handupres struct {
	Errcode int
}

type Handupnotice struct {
	Tableid int
	Chairid int
}

//Tableplayerleave
const (
	TableLeave_ToRoom = iota
	TableLeave_ChangeTable
)

type Tableplayerleavereq struct {
	Leavetype int //离开类型
}
type Tableplayerleaveres struct {
	Errcode   int
	Leavetype int //离开类型
}
type Tableplayerleavenotice struct {
	Userid  int64 //
	Tableid int
	Chairid int //
}

//
const (
	Kickout_toolong_handup = iota //入座举手超时
	Kickout_toolong_think         //对战考虑超时
	Kickout_room_closed           //房间已关闭
)

type Kickoutres struct {
	Errcode int
	Errmsg  string
}
