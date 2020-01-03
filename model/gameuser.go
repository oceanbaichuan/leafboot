package model

import "time"

type GameUserBackpack struct {
	UserID     int64 `gorm:"primary_key:true"` //用户ID
	PropID     int32
	PropType   int32
	PropCount  int32
	PropTime   time.Time
	UpdateTime time.Time
	CreateTime time.Time
}
type GameUserCurrency struct {
	UserID      int64 `gorm:"primary_key:true"` //用户ID
	GameCoin    int64
	PrizeTicket int32
	UpdateTime  time.Time
}

type GameUserDaybuysaleRecord struct {
	UserID       int64 `gorm:"primary_key:true"` //用户ID
	GameID       int32
	DaybuyTimes  int32
	DaybuyCoin   int32
	DaysaleTimes int32
	DaysaleCoins int32
	UpdateTime   time.Time
	DUpdateTime  time.Time
}
type GameUserDayProperty struct {
	UserID         int64 `gorm:"primary_key:true"` //用户ID
	OnlineTime     int32
	PlayTime       int32
	PlayCoin       int64
	PlayTax        int64
	GetPrizeticket int32
	HelpCoin       uint64
	HelpTimes      int32
	RechargeTimes  int32
	RechargeMoney  int32
	UpdateTime     time.Time
	DUpdateTime    time.Time
}

type GameUserOnline struct {
	UserID      int64 `gorm:"primary_key:true"` //用户ID
	AppID       string
	SiteID      int32
	GameID      int32
	ServerID    string
	ServerLevel int32
	CreateTime  time.Time
}

type GameUserProperty struct {
	UserID          int64 `gorm:"primary_key:true"` //用户ID
	OnlineTime      int32
	PlayTime        int32
	PlayCoin        int64
	PlayTax         int64
	PlayPrizeticket int32
	RechargeTimes   int32
	RechargeMoney   int32
	VipExp          int32
	GameExp         int32
	UpdateTime      time.Time
	DUpdateTime     time.Time
}
