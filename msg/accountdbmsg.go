package msg

import "time"

type PropInfo struct {
	Propid   int
	Proptype int
	Propnum  int
	Deadline time.Time
}
type UserPropChange struct {
	Propid   int
	Proptype int
	Propnum  int
	Deadline int64
}
type DBLoginRes struct {
	Userid                  int64
	Errcode                 int
	Useraccount             string    //用户账号
	Useraccountindex        int       //用户账号所属机器码的序号
	Useraccounttype         int       //账号类型
	Accountfreezed          bool      //账号是否冻结
	Usernickname            string    //用户昵称
	Usergamecoin            int64     //积分
	Userbankcoin            int64     //
	Usergoldbean            int       //金豆
	Allgetgoldbean          int       //
	Userregistechannelshort int       //注册短渠道
	Userregistetime         time.Time //注册时间
	Userregistesiteid       int
	Usergender              int8       //性别
	Usersysheadid           int        //系统头像id
	Userthirdheadurl        string     //第三方头像地址
	Userphonebinded         string     //绑定手机号
	VipExp                  int        //
	Gamewintimes            int        //
	Gamelosetimes           int        //
	Gameplaytime            int        //
	Platplaytime            int        //
	Gameonlinetime          int        //
	Platonlinetime          int        //
	Gamecoinplay            int64      //
	Gamecoinwin             int64      //
	Gamecoinlose            int64      //
	Othergameid             int        //
	Otherroomid             int        //
	Proplist                []PropInfo //道具列表
}

type Useronlinestatus struct {
	Userid      int64     `gorm:"primary_key:true"`
	Siteid      int       //
	Gameid      int       //
	Serverid    string    //
	Collectdate time.Time //
}
type Accountrdatareq struct {
	Userid   int64  //
	Gameid   int    //
	Serverid string //
	Loginip  string //
	Proplist []int  //道具列表
}
type Accountwdatareq struct {
	Userid         int64            //
	Gameid         int              //
	Gamecoin       int64            //增量
	Goldbean       int              //增量
	Proplist       []UserPropChange //道具列表增量
	Gametax        int64
	Gameplaytime   int   //
	Gameonlinetime int64 //
	Gamewintimes   int   //
	Gamelosetimes  int   //
	Gamecoinwin    int64 //
	Gamecoinlose   int64 //
}
