package model

import "time"

type GameDataRecord struct {
	UserID        int64 `gorm:"primary_key:true"` //用户ID
	GameID        int32 `gorm:"primary_key:true"` //用户ID
	HelpTimes     int32
	HelpCoin      int64
	RechargeTimes int32 //
	RechargeMoney int32
	PlayTime      int32
	PlayCoin      int64
	OnlineTime    int32
	Prizeticket   int32
	WinTimes      int32
	LoseTimes     int32
	GameExp       int32
	LoginTime     time.Time
	LoginIP       string
	UpdateTime    time.Time
}

func (a GameDataRecord) TableName() string {
	return "game_data_record"
}
