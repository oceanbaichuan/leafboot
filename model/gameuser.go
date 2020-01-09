package model

import "time"

type GameUserBackpack struct {
	UserID     int64     `gorm:"primary_key:true" json:"user_id"` //用户ID
	PropID     int32     `json:"prop_id"`
	PropType   int32     `json:"prop_type"`
	PropCount  int32     `json:"prop_count"`
	PropTime   int64     `json:"prop_time"`
	UpdateTime time.Time `json:"update_time"`
	CreateTime time.Time `json:"create_time"`
}

func (a GameUserBackpack) TableName() string {
	return "game_user_backpack"
}

type GameUserCurrency struct {
	UserID      int64     `gorm:"primary_key:true" json:"user_id"` //用户ID
	GameCoin    int64     `json:"game_coin"`
	PrizeTicket int32     `json:"prize_ticket"`
	UpdateTime  time.Time `json:"update_time"`
}

func (a GameUserCurrency) TableName() string {
	return "game_user_currency"
}

type GameUserDaybuysaleRecord struct {
	UserID       int64     `gorm:"primary_key:true" json:"user_id"` //用户ID
	GameID       int32     `json:"game_id"`
	DaybuyTimes  int32     `json:"daybuy_times"`
	DaybuyCoin   int32     `json:"daybuy_coin"`
	DaysaleTimes int32     `json:"daysale_times"`
	DaysaleCoins int32     `json:"daysale_coins"`
	UpdateTime   time.Time `json:"update_time"`
}

func (a GameUserDaybuysaleRecord) TableName() string {
	return "game_user_daybuysale_record"
}

type GameUserDayProperty struct {
	UserID          int64     `gorm:"primary_key:true" json:"user_id"` //用户ID
	OnlineTime      int32     `json:"online_time"`
	PlayTime        int32     `json:"play_time"`
	PlayCoin        int64     `json:"play_coin"`
	PlayTax         uint64    `json:"play_tax"`
	PlayPrizeticket int32     `json:"play_prizeticket"`
	HelpCoin        uint64    `json:"help_coin"`
	HelpTimes       int32     `json:"help_times"`
	RechargeTimes   int32     `json:"recharge_times"`
	RechargeMoney   int32     `json:"recharge_money"`
	UpdateTime      time.Time `json:"update_time"`
}

func (a GameUserDayProperty) TableName() string {
	return "game_user_dayproperty"
}

type GameUserOnline struct {
	UserID      int64     `gorm:"primary_key:true" json:"user_id"` //用户ID
	AppID       string    `json:"app_id"`
	SiteID      int32     `json:"site_id"`
	GameID      int32     `json:"game_id"`
	ServerID    string    `json:"server_id"`
	ServerLevel int32     `json:"server_level"`
	CreateTime  time.Time `json:"create_time"`
}

func (a GameUserOnline) TableName() string {
	return "game_user_online"
}

type GameUserProperty struct {
	UserID          int64     `gorm:"primary_key:true" json:"user_id"` //用户ID
	OnlineTime      int32     `json:"online_time"`
	PlayTime        int32     `json:"play_time"`
	PlayCoin        int64     `json:"play_coin"`
	PlayTax         uint64    `json:"play_tax"`
	PlayPrizeticket int32     `json:"play_prizeticket"`
	RechargeTimes   int32     `json:"recharge_times"`
	RechargeMoney   int32     `json:"recharge_money"`
	VipExp          int32     `json:"vip_exp"`
	GameExp         int32     `json:"game_exp"`
	UpdateTime      time.Time `json:"update_time"`
}

func (a GameUserProperty) TableName() string {
	return "game_user_property"
}
