package db

import (
	"github.com/hudgit2019/leafboot/model"
	"github.com/hudgit2019/leafboot/msg"
	"github.com/name5566/leaf/log"
)

const userdbName = "game_user_db"

func DeleteOnline(userid int64) {
	//Accountconn.dbconn.Table("QPGameUserDB.dbo.UserOnlineStatus").Where("userid = ?", userid).Delete()
	dbconn, err := GetDB(userdbName, userid)
	if err != nil {
		log.Error("DeleteOnline err:%v", err)
		return
	}
	dbconn.Model(&model.GameUserOnline{}).Delete(userid)
}
func SaveUserProperty(data *msg.Accountwdatareq) error {

	return nil
}
