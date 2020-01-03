package db

import (
	"github.com/hudgit2019/leafboot/model"
	"github.com/hudgit2019/leafboot/msg"
	"github.com/name5566/leaf/log"
)

const userdbName = "game_user_db"

func SelectUserInfo(datareq *msg.Accountrdatareq, account *msg.LoginRes) {
	dbconn, err := GetDB(userdbName, datareq.Userid)
	if err != nil {
		return
	}
	userdata := model.GameUserProperty{}
	dbconn.Find(&userdata, datareq.Userid)
	dbconn.Find(&userdata.GameUserCurrency, datareq.Userid)
	dbconn.Find(&userdata.GameUserBackpacks, datareq.Userid)
	dbconn.Find(&userdata.GameUserDaybuysaleRecord, datareq.Userid)
	dbconn.Find(&userdata.GameUserDayProperty, datareq.Userid)
	dbconn.Find(&userdata.GameUserOnline, datareq.Userid)
	account.GameCoin = userdata.GameUserCurrency.GameCoin
	account.GoldBean = userdata.GameUserCurrency.PrizeTicket
	account.AllGoldBean = userdata.PlayPrizeticket
	account.VipExp = userdata.VipExp
	account.PlatPlayTime = userdata.PlayTime
	account.PlatOnlineTime = userdata.OnlineTime
	account.OtherGameID = userdata.GameUserOnline.GameID
	account.OtherRoomID = userdata.GameUserOnline.ServerID
	for _, v := range userdata.GameUserBackpacks {
		prop := msg.PropInfo{
			Propid:   v.PropID,
			Proptype: v.PropType,
			Propnum:  v.PropCount,
			Proptime: v.PropTime,
		}
		account.Proplist = append(account.Proplist, prop)
	}
	log.Debug("SelectUserInfo %v", userdata)
}
func DeleteOnline(userid int64) {
	//Accountconn.dbconn.Table("QPGameUserDB.dbo.UserOnlineStatus").Where("userid = ?", userid).Delete()
	dbconn, err := GetDB(userdbName, userid)
	if err != nil {
		log.Error("DeleteOnline err:%v", err)
		return
	}
	dbconn.Delete(&model.GameUserOnline{UserID: userid})
}
func SaveUserProperty(data *msg.Accountwdatareq) error {

	return nil
}
