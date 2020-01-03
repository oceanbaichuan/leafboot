package model

import "time"

//AccountMainInfo 用户平台主信息
type AccountMainInfo struct {
	UserID                   int64  `gorm:"primary_key:true"` //用户ID
	Account                  string //
	Passwd                   string //
	AppID                    string //
	NickName                 string //
	SiteID                   int32  //
	ChannelID                int32  //
	AccountIndex             int32  //
	AccountType              int16
	AccountStatus            int8
	RegisteDate              time.Time
	RegisteIP                string
	RegisteMacid             string
	AccountSecurityInfo      AccountSecurityInfo
	AccountThirdpartInfo     AccountThirdpartInfo
	AccountLastloginInfo     AccountLastloginInfo
	GameUserBackpacks        []GameUserBackpack
	GameUserCurrency         GameUserCurrency
	GameUserDaybuysaleRecord GameUserDaybuysaleRecord
	GameUserDayProperty      GameUserDayProperty
	GameUserOnline           GameUserOnline
	GameUserProperty         GameUserProperty
}

//AccountSecurityInfo 用户安全信息
type AccountSecurityInfo struct {
	UserID     int64 `gorm:"primary_key:true"` //用户ID
	IdcID      string
	PhoneNum   string
	EmailAddr  string
	AlipayAddr string
	UpdateTime time.Time
}

//AccountThirdpartInfo 用户第三方信息
type AccountThirdpartInfo struct {
	UserID  int64 `gorm:"primary_key:true"` //用户ID
	OpenID  string
	HeadUrl string
}

//AccountLastloginInfo 用户最近登录平台信息
type AccountLastloginInfo struct {
	UserID          int64 `gorm:"primary_key:true"` //用户ID
	LastLoginIP     string
	LastLoginMacid  string
	LastLoginAppid  string
	LastLoginChanid int32
	LastLoginDevid  int8
	LoginTime       time.Time
}
