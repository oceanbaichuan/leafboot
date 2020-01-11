package db

import (
	"fmt"
	"time"

	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/model"
	"github.com/hudgit2019/leafboot/msg"
	myredis "github.com/hudgit2019/leafboot/redis"
	"github.com/jinzhu/gorm"
	"github.com/name5566/leaf/log"
)

const userdbName = "game_user_db"

func SelectUserInfo(datareq *msg.Accountrdatareq, account *msg.LoginRes) {

	//读取游戏信息
	dbrconn, err := GetDB(userdbName, datareq.Userid, DBFLAG_R, true)
	if err != nil {
		log.Debug("dbuser SelectUserInfo account:%s, err:%v", datareq.Account, err)
		return
	}
	dbwconn, err := GetDB(userdbName, datareq.Userid, DBFLAG_RW, true)
	if err != nil {
		log.Debug("dbuser SelectUserInfo account:%s, err:%v", datareq.Account, err)
		return
	}
	userdata := model.GameUserProperty{}
	var propList []msg.PropInfo
	userCurrency := model.GameUserCurrency{}
	//daybuysaleRecord := model.GameUserDaybuysaleRecord{}
	userDayProperty := model.GameUserDayProperty{}
	userOnline := model.GameUserOnline{}
	if result, err := myredis.SelectUserInfo(userdbName, datareq); err == nil {
		userCurrency = result[0].(model.GameUserCurrency)
		userdata = result[1].(model.GameUserProperty)
		userDayProperty = result[2].(model.GameUserDayProperty)
		propList = result[3].([]msg.PropInfo)
		userOnline = result[4].(model.GameUserOnline)
		userdata.UserID = datareq.Userid
		userDayProperty.UserID = datareq.Userid
		userCurrency.UserID = datareq.Userid
		userOnline.UserID = datareq.Userid
		//跨天
		if base.CheckDateDiff(userDayProperty.UpdateTime, time.Now(), "day", 1) {
			userDayProperty = model.GameUserDayProperty{
				UserID:     datareq.Userid,
				UpdateTime: time.Now(),
			}
			myredis.WriteUserDayProperty(userdbName, userDayProperty)
		}
		//更新redis有效时长
		var propIDKeys []string
		var propTTFs []time.Duration
		for _, v := range propList {
			propIDKeys = append(propIDKeys, fmt.Sprintf(myredis.USER_PROPID_SET_KEY, datareq.Userid, v.Propid))
			propTTFs = append(propTTFs, myredis.REDIS_ACCOUNT_TTF)
			propIDKeys = append(propIDKeys, fmt.Sprintf(myredis.USER_PROPINFO_KEY, datareq.Userid, v.Propid, v.Proptype))
			propTTFs = append(propTTFs, myredis.REDIS_ACCOUNT_TTF)
		}
		ttfKeys := []string{
			fmt.Sprintf(myredis.USER_CURRENCY_KEY, datareq.Userid),
			fmt.Sprintf(myredis.USER_PROPERTY_KEY, datareq.Userid),
			fmt.Sprintf(myredis.USER_DAYPROPERTY_KEY, datareq.Userid),
			fmt.Sprintf(myredis.USER_BACKPACK_SET_KEY, datareq.Userid),
		}
		ttfValues := []time.Duration{
			myredis.REDIS_ACCOUNT_TTF,
			myredis.REDIS_ACCOUNT_TTF,
			myredis.REDIS_ACCOUNT_TTF,
			myredis.REDIS_ACCOUNT_TTF,
		}
		ttfValues = append(ttfValues, propTTFs...)
		ttfKeys = append(ttfKeys, propIDKeys...)
		UpdateUserDataTTF(ttfKeys, datareq.Userid, ttfValues)
		//非断线重连，插入在线标记
		if userOnline.GameID == 0 &&
			datareq.Gameid > 0 &&
			datareq.Serverid != "" {
			onlineInfo := model.GameUserOnline{
				UserID:      datareq.Userid,
				GameID:      datareq.Gameid,
				ServerID:    datareq.Serverid,
				AppID:       account.AppID,
				SiteID:      account.SiteID,
				ServerLevel: datareq.ServerLevel,
				CreateTime:  time.Now(),
			}
			dbwconn.Save(&onlineInfo)
			myredis.WriteUserOnline(userdbName, onlineInfo)
		}
	} else {
		dbrconn.Find(&userdata, datareq.Userid)
		var backpacks []model.GameUserBackpack
		//daybuysaleRecord := model.GameUserDaybuysaleRecord{}
		dbrconn.Find(&userCurrency, datareq.Userid)
		dbrconn.Find(&backpacks, datareq.Userid)
		//dbconn.Find(&daybuysaleRecord, datareq.Userid)
		dbrconn.Find(&userDayProperty, datareq.Userid)
		dbrconn.Find(&userOnline, datareq.Userid)
		//道具
		for _, v := range backpacks {
			prop := msg.PropInfo{
				Propid:   v.PropID,
				Proptype: v.PropType,
				Propnum:  v.PropCount,
				Proptime: v.PropTime,
			}
			propList = append(propList, prop)
		}
		//数据写入redis
		onlineInfo := model.GameUserOnline{
			UserID:      datareq.Userid,
			GameID:      datareq.Gameid,
			ServerID:    datareq.Serverid,
			AppID:       account.AppID,
			SiteID:      account.SiteID,
			ServerLevel: datareq.ServerLevel,
			CreateTime:  time.Now(),
		}
		if userOnline.GameID == 0 &&
			datareq.Gameid > 0 &&
			datareq.Serverid != "" {
			dbwconn.Save(&onlineInfo)
		}
		userdata.UserID = datareq.Userid
		userDayProperty.UserID = datareq.Userid
		userCurrency.UserID = datareq.Userid
		userOnline.UserID = datareq.Userid
		//跨天
		if base.CheckDateDiff(userDayProperty.UpdateTime, time.Now(), "day", 1) {
			userDayProperty = model.GameUserDayProperty{
				UserID:     datareq.Userid,
				UpdateTime: time.Now(),
			}
		}
		log.Debug("backpacks :%v", backpacks)
		myredis.WriteUserInfoRedis(userdbName, userdata, userDayProperty, userCurrency, backpacks, onlineInfo)
	}
	//货币
	account.GameCoin = userCurrency.GameCoin
	account.GoldBean = userCurrency.PrizeTicket
	//平台属性
	account.AllGoldBean = userdata.PlayPrizeticket
	account.VipExp = userdata.VipExp
	account.PlatPlayTime = userdata.PlayTime
	account.PlatOnlineTime = userdata.OnlineTime
	account.AllGameExp = userdata.GameExp
	account.PlatRechargeTimes = userdata.RechargeTimes
	account.PlatRechargeMoney = userdata.RechargeMoney
	account.PlatTax = userdata.PlayTax
	account.PlatPrizeTicket = userdata.PlayPrizeticket //平台总获得奖券
	//每日属性
	account.PlatDayPlayTime = userDayProperty.PlayTime
	account.PlatDayOnlineTime = userDayProperty.OnlineTime
	account.PlatDayRechargeTimes = userDayProperty.RechargeTimes
	account.PlatDayRechargeMoney = userDayProperty.RechargeMoney
	account.PlatDayTax = userDayProperty.PlayTax
	account.PlatDayPrizeTicket = userDayProperty.PlayPrizeticket //平台总获得奖券
	account.OtherGameID = userOnline.GameID
	account.OtherRoomID = userOnline.ServerID
	account.Proplist = propList
	log.Debug("SelectUserInfo userdata: %v", userdata)
	log.Debug("SelectUserInfo userDayProperty: %v", userDayProperty)
	log.Debug("SelectUserInfo userOnline: %v", userOnline)
	log.Debug("SelectUserInfo userCurrency: %v", userCurrency)
	log.Debug("SelectUserInfo propList: %v", propList)
}
func DeleteOnline(userid int64) {
	//Accountconn.dbconn.Table("QPGameUserDB.dbo.UserOnlineStatus").Where("userid = ?", userid).Delete()
	dbwconn, err := GetDB(userdbName, userid, DBFLAG_RW, true)
	if err != nil {
		log.Error("DeleteOnline userid:%v err:%v", userid, err)
		return
	}
	dbwconn.Delete(&model.GameUserOnline{UserID: userid})
	myredis.DeleteOnline(userdbName, userid)
}
func SaveUserProperty(data msg.Accountwdatareq) error {
	dbwconn, err := GetDB(userdbName, data.UserID, DBFLAG_RW, true)
	if err != nil {
		log.Error("SaveUserProperty userid:%v err:%v", data.UserID, err)
		return err
	}
	//更新货币
	if data.GameCoin != 0 || data.PrizeTicket != 0 {
		dbrow := dbwconn.Model(&model.GameUserCurrency{UserID: data.UserID}).Updates(map[string]interface{}{
			"game_coin":    gorm.Expr("game_coin + ?", data.GameCoin),
			"prize_ticket": gorm.Expr("prize_ticket + ?", data.PrizeTicket),
		})
		if dbrow.RowsAffected <= 0 {
			dbwconn.Save(&model.GameUserCurrency{
				UserID:      data.UserID,
				GameCoin:    data.GameCoin,
				PrizeTicket: data.PrizeTicket,
				UpdateTime:  time.Now(),
			})
		}
	}
	//更新总属性
	dbrow := dbwconn.Model(&model.GameUserProperty{UserID: data.UserID}).Updates(map[string]interface{}{
		"online_time":      gorm.Expr("online_time + ?", data.GameOnlinetime),
		"play_time":        gorm.Expr("play_time + ?", data.GamePlaytime),
		"play_coin":        gorm.Expr("play_coin + ?", data.GameCoin),
		"play_prizeticket": gorm.Expr("play_prizeticket + ?", data.PrizeTicket),
		"game_exp":         gorm.Expr("game_exp + ?", data.GameExp),
		"play_tax":         gorm.Expr("play_tax + ?", data.GameTax),
	})
	if dbrow.RowsAffected <= 0 {
		dbwconn.Save(&model.GameUserProperty{
			UserID:          data.UserID,
			OnlineTime:      data.GameOnlinetime,
			PlayTime:        data.GamePlaytime,
			PlayCoin:        data.GameCoin,
			PlayPrizeticket: data.PrizeTicket,
			GameExp:         data.GameExp,
			PlayTax:         data.GameTax,
			UpdateTime:      time.Now(),
		})
	}
	//更新每日属性
	dbrow = dbwconn.Model(&model.GameUserDayProperty{UserID: data.UserID}).Updates(map[string]interface{}{
		"online_time":      gorm.Expr("online_time + ?", data.GameOnlinetime),
		"play_time":        gorm.Expr("play_time + ?", data.GamePlaytime),
		"play_coin":        gorm.Expr("play_coin + ?", data.GameCoin),
		"play_prizeticket": gorm.Expr("play_prizeticket + ?", data.PrizeTicket),
		"game_exp":         gorm.Expr("game_exp + ?", data.GameExp),
		"play_tax":         gorm.Expr("play_tax + ?", data.GameTax),
		"help_times":       gorm.Expr("help_times + ?", data.HelpTimes),
	})
	if dbrow.RowsAffected <= 0 {
		dbwconn.Save(&model.GameUserDayProperty{
			UserID:          data.UserID,
			OnlineTime:      data.GameOnlinetime,
			PlayTime:        data.GamePlaytime,
			PlayCoin:        data.GameCoin,
			PlayPrizeticket: data.PrizeTicket,
			PlayTax:         data.GameTax,
			UpdateTime:      time.Now(),
		})
	}
	//更新道具
	for _, v := range data.PropList {
		prop := &model.GameUserBackpack{
			UserID:   data.UserID,
			PropID:   v.Propid,
			PropType: v.Proptype,
		}
		dbrow = dbwconn.Model(prop).Updates(map[string]interface{}{
			"prop_count": gorm.Expr("prop_count + ?", v.Propnum),
			"prop_time":  gorm.Expr("prop_time + ?", v.Proptime),
		})
		if dbrow.RowsAffected <= 0 {
			dbwconn.Save(&model.GameUserBackpack{
				UserID:     data.UserID,
				PropID:     v.Propid,
				PropType:   v.Proptype,
				PropCount:  v.Propnum,
				PropTime:   v.Proptime,
				UpdateTime: time.Now(),
			})
		}
	}
	return myredis.SaveUserProperty(userdbName, data)
}
func UpdateUserDataTTF(keys []string, userID int64, ttfs []time.Duration) {
	myredis.UpdateRedisTTF(userdbName, userID, keys, ttfs)
}
