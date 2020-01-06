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

//LeaveSceneReq 进入业务服
type LeaveSceneReq struct {
	UserID int64 //用户ID
}
type LeaveSceneRes struct {
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
	Errcode        int32
	Account        string    //用户账号
	AcIndex        int32     //用户账号所属机器码的序号
	AcType         int16     //账号类型
	AppID          string    //对外产品ID
	SiteID         int32     //对内产品ID
	Freezed        int8      //账号是否冻结
	NickName       string    //用户昵称
	GameCoin       int64     //积分
	BankCoin       int64     //
	GoldBean       int32     //金豆
	AllGoldBean    int32     //
	RegChan        int32     //注册短渠道
	RegTime        time.Time //注册时间
	RegSiteID      int32
	Gender         int8   //性别
	HeadID         int32  //系统头像id
	ThirdHeadUrl   string //第三方头像地址
	Phonebinded    string //绑定手机号
	VipExp         int32  //
	GameExp        int32  //
	GameWinTimes   int32  //
	GameLoseTimes  int32  //
	GamePlayTime   int32  //
	PlatPlayTime   int32  //
	GameOnlineTime int32  //
	PlatOnlineTime int32  //
	GameCoinPlay   int64  //
	GameCoinWin    int64  //
	GameCoinLose   int64  //
	OtherGameID    int32  //
	OtherRoomID    string //
	GameStatus     int8
	Proplist       []PropInfo //道具列表
}

type TablePlayer struct {
	Userid        int64 //
	Siteid        int32
	Nickname      string //
	Gamecoin      int64
	Bankcoin      int64
	Goldbean      int32
	Sysheadid     int32
	Thirdheadurl  string //
	Gender        int8
	Vipexp        int32 //
	Gamewintimes  int32 //
	Gamelosetimes int32 //
	Tableid       int32
	Chairid       int32 //
	Gamestatus    int32 //
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
	Tableid int32 //
	Chairid int32 //
}

type Sitdownres struct {
	Errcode int32 //错误码
	Tableid int32 //桌子号
	Chairid int32 //椅子号
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
	Errcode int32
}

type Handupnotice struct {
	Tableid int32
	Chairid int32
}

//Tableplayerleave
const (
	TableLeave_ToRoom = iota
	TableLeave_ChangeTable
)

type Tableplayerleavereq struct {
	Leavetype int32 //离开类型
}
type Tableplayerleaveres struct {
	Errcode   int32
	Leavetype int32 //离开类型
}
type Tableplayerleavenotice struct {
	Userid  int64 //
	Tableid int32
	Chairid int32 //
}

//
const (
	Kickout_toolong_handup = iota //入座举手超时
	Kickout_toolong_think         //对战考虑超时
	Kickout_room_closed           //房间已关闭
)

type Kickoutres struct {
	Errcode int32
	Errmsg  string
}
