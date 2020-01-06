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
type DBLoginRes struct {
	Userid                  int64
	Errcode                 int32
	Useraccount             string    //用户账号
	Useraccountindex        int32     //用户账号所属机器码的序号
	Useraccounttype         int16     //账号类型
	Accountfreezed          int8      //账号是否冻结
	Usernickname            string    //用户昵称
	Usergamecoin            int64     //积分
	Userbankcoin            int64     //
	Usergoldbean            int32     //金豆
	Allgetgoldbean          int32     //
	Userregistechannelshort int32     //注册短渠道
	Userregistetime         time.Time //注册时间
	Userregistesiteid       int32
	Usergender              int8       //性别
	Usersysheadid           int32      //系统头像id
	Userthirdheadurl        string     //第三方头像地址
	Userphonebinded         string     //绑定手机号
	VipExp                  int32      //
	Gamewintimes            int32      //
	Gamelosetimes           int32      //
	Gameplaytime            int32      //
	Platplaytime            int32      //
	Gameonlinetime          int32      //
	Platonlinetime          int32      //
	Gamecoinplay            int64      //
	Gamecoinwin             int64      //
	Gamecoinlose            int64      //
	Othergameid             int32      //
	Otherroomid             int32      //
	Proplist                []PropInfo //道具列表
}

type Useronlinestatus struct {
	Userid      int64     `gorm:"primary_key:true"`
	Siteid      int32     //
	Gameid      int32     //
	Serverid    string    //
	Collectdate time.Time //
}
type Accountrdatareq struct {
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
	PropList       []UserPropChange //道具列表增量
	GameTax        int64
	GamePlaytime   int32 //
	GameOnlinetime int64 //
	GameWintimes   int32 //
	GameLosetimes  int32 //
	GameCoinwin    int64 //
	GameCoinlose   int64 //
}
