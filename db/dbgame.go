package db

import (
	"time"

	"github.com/hudgit2019/leafboot/model"
	"github.com/hudgit2019/leafboot/msg"
	myredis "github.com/hudgit2019/leafboot/redis"
	"github.com/jinzhu/gorm"
	"github.com/name5566/leaf/log"
)

const gamedbName = "game_data_db"

func SelectGameData(datareq *msg.Accountrdatareq, account *msg.LoginRes) {

	//读取游戏信息
	dbrconn, err := GetDB(gamedbName, datareq.Userid, DBFLAG_R, true)
	if err != nil {
		log.Debug("dbgame SelectGameData account:%s, err:%v", datareq.Account, err)
		return
	}
	dbwconn, err := GetDB(gamedbName, datareq.Userid, DBFLAG_RW, true)
	if err != nil {
		log.Debug("dbgame SelectGameData account:%s, err:%v", datareq.Account, err)
		return
	}
	gamedata := model.GameDataRecord{}
	if result, err := myredis.SelectGameData(gamedbName, datareq); err != nil {
		dbRow := dbrconn.Where("user_id = ? and game_id = ?", datareq.Userid, datareq.Gameid).Find(&gamedata)
		log.Debug("SelectGameData %v", gamedata)
		if dbRow.RowsAffected > 0 {
			dbwconn.Model(&model.GameDataRecord{UserID: datareq.Userid,
				GameID: datareq.Gameid}).Update(model.GameDataRecord{
				LoginIP:   datareq.Loginip,
				LoginTime: time.Now(),
			})
		} else {
			dbwconn.Save(&model.GameDataRecord{
				UserID:     datareq.Userid,
				GameID:     datareq.Gameid,
				LoginIP:    datareq.Loginip,
				LoginTime:  time.Now(),
				UpdateTime: time.Now(),
			})
		}
	} else {
		gamedata = result[0].(model.GameDataRecord)
		dbwconn.Save(&model.GameDataRecord{
			UserID:     datareq.Userid,
			GameID:     datareq.Gameid,
			LoginIP:    datareq.Loginip,
			LoginTime:  time.Now(),
			UpdateTime: time.Now(),
		})
	}
	account.GameExp = gamedata.GameExp
	account.GameWinTimes = gamedata.WinTimes
	account.GameLoseTimes = gamedata.LoseTimes
	account.GamePlayTime = gamedata.PlayTime
	account.GameOnlineTime = gamedata.OnlineTime
	account.GameCoinPlay = gamedata.PlayCoin
}
func SaveGameData(data msg.Accountwdatareq) error {
	//更新游戏属性
	dbwconn, err := GetDB(gamedbName, data.UserID, DBFLAG_RW, true)
	if err != nil {
		return err
	}
	dbrow := dbwconn.Model(&model.GameDataRecord{UserID: data.UserID, GameID: data.GameID}).Updates(map[string]interface{}{
		"online_time": gorm.Expr("online_time + ?", data.GameOnlinetime),
		"play_time":   gorm.Expr("play_time + ?", data.GamePlaytime),
		"play_coin":   gorm.Expr("play_coin + ?", data.GameCoin),
		"prizeticket": gorm.Expr("prizeticket + ?", data.PrizeTicket),
		"game_exp":    gorm.Expr("game_exp + ?", data.GameExp),
		"help_times":  gorm.Expr("help_times + ?", data.HelpTimes),
		"help_coin":   gorm.Expr("help_coin + ?", data.HelpTimes),
		"win_times":   gorm.Expr("win_times + ?", data.GameWintimes),
		"lose_times":  gorm.Expr("lose_times + ?", data.GameLosetimes),
	})
	if dbrow.RowsAffected <= 0 {
		dbwconn.Save(&model.GameDataRecord{
			UserID:      data.UserID,
			GameID:      data.GameID,
			OnlineTime:  data.GameOnlinetime,
			PlayTime:    data.GamePlaytime,
			PlayCoin:    data.GameCoin,
			Prizeticket: data.PrizeTicket,
			GameExp:     data.GameExp,
			HelpTimes:   data.HelpTimes,
			WinTimes:    data.GameWintimes,
			LoseTimes:   data.GameLosetimes,
			LoginTime:   time.Now(),
			UpdateTime:  time.Now(),
		})
	}
	return myredis.SaveGameData(gamedbName, data)
}
