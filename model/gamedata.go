package model

import "time"

type GameDataRecord struct {
	UserID        int64     `gorm:"primary_key:true" json:"user_id"` //用户ID
	GameID        int32     `gorm:"primary_key:true" json:"game_id"` //用户ID
	HelpTimes     int32     `json:"help_times"`
	HelpCoin      int64     `json:"help_coin"`
	RechargeTimes int32     `json:"recharge_times"`
	RechargeMoney int32     `json:"recharge_money"`
	PlayTime      int32     `json:"play_time"`
	PlayCoin      int64     `json:"play_coin"`
	OnlineTime    int32     `json:"online_time"`
	Prizeticket   int32     `json:"prizeticket"`
	WinTimes      int32     `json:"win_times"`
	LoseTimes     int32     `json:"lose_times"`
	GameExp       int32     `json:"game_exp"`
	LoginTime     time.Time `json:"login_time"`
	LoginIP       string    `json:"login_ip"`
	UpdateTime    time.Time `json:"update_time"`
}

func (a GameDataRecord) TableName() string {
	return "game_data_record"
}
