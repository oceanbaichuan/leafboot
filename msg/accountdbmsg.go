package msg

import "time"

type PropInfo struct {
	Propid   int32
	Proptype int32
	Propnum  int32
	Proptime int64
}
type UserPropChange struct {
	Propid   int32
	Proptype int32
	Propnum  int32
	Proptime int64
}
type LoginRes struct {
	Userid       int64
	Errcode      int32
	Account      string //用户账号
	Token        string
	AcIndex      int32     //用户账号所属机器码的序号
	AcType       int16     //账号类型
	AppID        string    //对外产品ID
	SiteID       int32     //对内产品ID
	IsValid      int8      //账号是否冻结
	NickName     string    //用户昵称
	GameCoin     int64     //积分
	BankCoin     int64     //
	GoldBean     int32     //金豆
	AllGoldBean  int32     //
	RegChan      int32     //注册短渠道
	RegTime      time.Time //注册时间
	RegSiteID    int32     //注册站点
	Gender       int8      //性别
	HeadID       int32     //系统头像id
	ThirdHeadUrl string    //第三方头像地址
	Phonebinded  string    //绑定手机号
	//平台属性
	VipExp            int32 //
	AllGameExp        int32 //
	PlatPlayTime      int32 //
	PlatOnlineTime    int32 //
	PlatPlayCoin      int64 //
	PlatPrizeTicket   int32 //平台总获得奖券
	PlatRechargeTimes int32
	PlatRechargeMoney int32  //分为单位
	PlatTax           uint64 //
	//每日平台属性
	PlatDayPlayTime      int32 //
	PlatDayOnlineTime    int32 //
	PlatDayPlayCoin      int64 //
	PlatDayPrizeTicket   int32 //平台总获得奖券
	PlatDayRechargeTimes int32
	PlatDayRechargeMoney int32  //分为单位
	PlatDayTax           uint64 //
	PlatHelpTimes        int32  //
	//游戏属性
	GameExp        int32  //
	GameWinTimes   int32  //
	GameLoseTimes  int32  //
	GamePlayTime   int32  //
	GameOnlineTime int32  //
	GameCoinPlay   int64  //
	GameCoinWin    int64  //
	GameCoinLose   int64  //
	OtherGameID    int32  //
	OtherRoomID    string //
	GameStatus     int8
	Proplist       []PropInfo //道具列表
}

type Useronlinestatus struct {
	Userid      int64     `gorm:"primary_key:true"`
	Siteid      int32     //
	Gameid      int32     //
	Serverid    string    //
	Collectdate time.Time //
}
type Accountrdatareq struct {
	Account     string
	Passwd      string
	Token       string  //登录token
	LoginRoute  string  //登录来源，Login,Lobby,DZMJ
	Userid      int64   //
	Gameid      int32   //
	Serverid    string  //
	ServerLevel int32   //
	Loginip     string  //
	MacID       string  //
	ChannelID   int32   //
	AppID       string  //
	DevType     int8    //
	Proplist    []int32 //道具列表
}
type Accountwdatareq struct {
	UserID         int64            //
	GameID         int32            //
	GameCoin       int64            //增量
	PrizeTicket    int32            //增量
	HelpTimes      int32            //增量
	PropList       []UserPropChange //道具列表增量
	GameTax        uint64
	GameExp        int32
	GamePlaytime   int32 //
	GameOnlinetime int32 //
	GameWintimes   int32 //
	GameLosetimes  int32 //
	GameCoinwin    int64 //
	GameCoinlose   int64 //
}
