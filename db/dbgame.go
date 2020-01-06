package db

import (
	"time"

	"github.com/hudgit2019/leafboot/model"
	"github.com/hudgit2019/leafboot/msg"
	"github.com/name5566/leaf/log"
)

const gamedbName = "game_data_db"

func SelectGameData(datareq *msg.Accountrdatareq, account *msg.LoginRes) {
	dbconn, err := GetDB(gamedbName, datareq.Userid)
	if err != nil {
		return
	}
	gamedata := model.GameDataRecord{}
	dbRow := dbconn.Where("user_id = ? and game_id = ?", datareq.Userid, datareq.Gameid).Find(&gamedata)
	log.Debug("SelectGameData %v", gamedata)
	account.GameExp = gamedata.GameExp
	account.GameWinTimes = gamedata.WinTimes
	account.GameLoseTimes = gamedata.LoseTimes
	account.GamePlayTime = gamedata.PlayTime
	account.GameOnlineTime = gamedata.OnlineTime
	account.GameCoinPlay = gamedata.PlayCoin
	if dbRow.RowsAffected > 0 {
		dbconn.Model(&model.GameDataRecord{UserID: datareq.Userid,
			GameID: datareq.Gameid}).Update(model.GameDataRecord{
			LoginIP:   datareq.Loginip,
			LoginTime: time.Now(),
		})
	} else {
		dbconn.Save(&model.GameDataRecord{
			UserID:     datareq.Userid,
			GameID:     datareq.Gameid,
			LoginIP:    datareq.Loginip,
			LoginTime:  time.Now(),
			UpdateTime: time.Now(),
		})
	}
}
