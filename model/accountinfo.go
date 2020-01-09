package model

import "time"

//AccountMainInfo 用户平台主信息

type AccountMainInfo struct {
	UserID        int64     `gorm:"primary_key:true" json:"user_id"` //用户ID
	Account       string    `json:"account"`                         //
	Passwd        string    `json:"passwd"`                          //
	AppID         string    `json:"app_id"`                          //
	NickName      string    `json:"nick_name"`                       //
	SiteID        int32     `json:"site_id"`                         //
	Gender        int8      `json:"gender"`
	ChannelID     int32     `json:"channel_id"`    //
	AccountIndex  int32     `json:"account_index"` //
	AccountType   int16     `json:"account_type"`
	AccountStatus int8      `json:"account_status"`
	RegisteDate   time.Time `json:"registe_date"`
	RegisteIP     string    `json:"registe_ip"`
	RegisteMacid  string    `json:"registe_macid"`
}

func (a AccountMainInfo) TableName() string {
	return "account_main_info"
}

//AccountSecurityInfo 用户安全信息
type AccountSecurityInfo struct {
	UserID     int64     `gorm:"primary_key:true" json:"user_id"` //用户ID
	IdcID      string    `json:"idc_id"`
	PhoneNum   string    `json:"phone_num"`
	EmailAddr  string    `json:"email_addr"`
	AlipayAddr string    `json:"alipay_addr"`
	UpdateTime time.Time `json:"update_time"`
}

func (a AccountSecurityInfo) TableName() string {
	return "account_security_info"
}

//AccountThirdpartInfo 用户第三方信息
type AccountThirdpartInfo struct {
	UserID  int64  `gorm:"primary_key:true" json:"user_id"` //用户ID
	OpenID  string `json:"open_id"`
	HeadUrl string `json:"head_url"`
}

func (a AccountThirdpartInfo) TableName() string {
	return "account_thirdpart_info"
}

//AccountLastloginInfo 用户最近登录平台信息
type AccountLastloginInfo struct {
	UserID          int64     `gorm:"primary_key:true" json:"user_id"` //用户ID
	LastLoginIP     string    `json:"last_login_ip"`
	LastLoginMacid  string    `json:"last_login_macid"`
	LastLoginAppid  string    `json:"last_login_appid"`
	LastLoginChanid int32     `json:"last_login_chanid"`
	LastLoginDevid  int8      `json:"last_login_devid"`
	LoginTime       time.Time `json:"login_time"`
}

func (a AccountLastloginInfo) TableName() string {
	return "account_lastlogin_info"
}
