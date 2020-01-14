package msg

import "time"

//EnterScenceReq 进入业务服,各业务可自由筛选字段但结构名不能变
type EnterScenceReq struct {
	Account     string //账号
	Passwd      string //密码
	Token       string //token
	UserID      int64  //用户ID
	AppID       string //产品ID
	MacID       string //机器码
	OsType      int8   //1:pc 2:Android 3:IOS
	LoginIP     string //登录IP
	OpenID      string //第三方OpenID
	UnionID     string //第三方UnionID
	LoginMethod int8   //第三方登录方式 1，2,3……
	ChannelID   int32  //短渠道号
	LongChannel string //长渠道号
	BigVer      string //大版本号
	SmallVer    string //小版本号
	UserSource  string //用户来源,例如同1个渠道不同的位置
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
type ClientUserRes struct {
	UserID     int64
	ErrCode    int32
	Account    string //用户账号
	Token      string
	AcIndex    int32     //用户账号所属机器码的序号
	AcType     int16     //账号类型
	AppID      string    //对外产品ID
	SID        int32     //对内产品ID
	IsValid    int8      //账号是否冻结
	NickName   string    //用户昵称
	GameCoin   int64     //积分
	PTicket    int32     //金豆
	AllPTicket int32     //
	RegChan    int32     //注册短渠道
	RegTime    time.Time //注册时间
	RegSID     int32
	Gender     int8   //性别
	HeadID     int32  //系统头像id
	THUrl      string //第三方头像地址
	PhNum      string //绑定手机号
	GameID     int32
	ServerID   string
	PropList   []PropInfo //道具列表
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

type LoginAgainRes struct {
	UserID       int64
	ErrCode      int32
	Account      string //用户账号
	Token        string
	AcIndex      int32     //用户账号所属机器码的序号
	AcType       int16     //账号类型
	AppID        string    //对外产品ID
	SID          int32     //对内产品ID
	IsValid      int8      //账号是否冻结
	NickName     string    //用户昵称
	GameCoin     int64     //积分
	PTicket      int32     //金豆
	AllPTicket   int32     //
	RegChan      int32     //注册短渠道
	RegTime      time.Time //注册时间
	RegSID       int32
	Gender       int8   //性别
	HeadID       int32  //系统头像id
	THUrl        string //第三方头像地址
	PhNum        string //绑定手机号
	GameID       int32
	ServerID     string
	PropList     []PropInfo //道具列表
	Tableplayers []TablePlayer
}

//入座
const (
	SitErr_NoFitLocation = iota
	SitErr_HaveSit
	SitErr_ForbidWatch
	SitErr_NoLogin
)

type SitDownReq struct {
	Levelid int32 //桌子等级
	Tableid int32 //
	Chairid int32 //
}

type SitDownRes struct {
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

type HandUpReq struct {
	Handup int8
}

type HandUpRes struct {
	Errcode int32
}

type HandUpNotice struct {
	Tableid int32
	Chairid int32
}

//Tableplayerleave
const (
	TableLeave_ToRoom = iota
	TableLeave_ChangeTable
)

type TablePlayerLeaveReq struct {
	Leavetype int32 //离开类型
}
type TablePlayerLeaveRes struct {
	Errcode   int32
	Leavetype int32 //离开类型
}
type TablePlayerLeaveNotice struct {
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

type KickOutRes struct {
	Errcode int32
	Errmsg  string
}
