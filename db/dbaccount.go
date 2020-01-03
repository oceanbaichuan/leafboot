package db

import (
	"github.com/hudgit2019/leafboot/msg"

	"github.com/hudgit2019/leafboot/model"
	"github.com/name5566/leaf/log"
)

const accountdbName = "plat_account_db"

func intit() {
}
func SelectAccount(datareq *msg.Accountrdatareq, account *msg.LoginRes) error {
	//读取账户基本信息
	accountInfo := model.AccountMainInfo{}
	dbconn, err := GetDB(accountdbName, datareq.Userid)
	if err != nil {
		return err
	}
	row := dbconn.Find(&accountInfo, datareq.Userid)

	if row.RowsAffected < 1 {
		account.Errcode = msg.LoginErr_nouser
		return row.Error
	}
	if accountInfo.AccountStatus == 0 {
		account.Errcode = msg.LoginErr_accountforbidden
		return nil
	}
	dbconn.Find(&accountInfo.AccountSecurityInfo, datareq.Userid)
	dbconn.Find(&accountInfo.AccountThirdpartInfo, datareq.Userid)
	dbconn.Find(&accountInfo.AccountLastloginInfo, datareq.Userid)
	account.Userid = accountInfo.UserID
	account.Account = accountInfo.Account
	account.AcIndex = accountInfo.AccountIndex
	account.AcType = accountInfo.AccountType
	account.Freezed = accountInfo.AccountStatus
	account.NickName = accountInfo.NickName
	account.RegChan = accountInfo.ChannelID
	account.RegTime = accountInfo.RegisteDate
	account.RegSiteID = accountInfo.SiteID
	account.Gender = accountInfo.Gender
	account.ThirdHeadUrl = accountInfo.AccountThirdpartInfo.HeadUrl //第三方头像地址
	account.Phonebinded = accountInfo.AccountSecurityInfo.PhoneNum  //绑定手机号
	log.Debug("%v", accountInfo)
	//tx.Commit()
	return nil
}
