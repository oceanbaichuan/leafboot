package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/model"
	"github.com/hudgit2019/leafboot/msg"
	myredis "github.com/hudgit2019/leafboot/redis"
	"github.com/name5566/leaf/log"
)

const accountdbName = "plat_account_db"

func intit() {
}
func SelectAccount(datareq *msg.Accountrdatareq, account *msg.LoginRes) error {
	//读取账户基本信息
	dbrconn, err := GetDB(accountdbName, datareq.Userid, DBFLAG_R, true)
	if err != nil {
		log.Debug("dbaccount SelectAccount account:%s, err:%v", datareq.Account, err)
		return err
	}
	dbwconn, err := GetDB(accountdbName, datareq.Userid, DBFLAG_RW, true)
	if err != nil {
		log.Debug("dbaccount SelectAccount account:%s, err:%v", datareq.Account, err)
		return err
	}

	loginInfo := model.AccountLastloginInfo{
		UserID:          datareq.Userid,
		LastLoginIP:     datareq.Loginip,
		LastLoginMacid:  datareq.MacID,
		LastLoginAppid:  datareq.AppID,
		LastLoginChanid: datareq.ChannelID,
		LastLoginDevid:  datareq.DevType,
		LoginTime:       time.Now(),
	}
	accountInfo := model.AccountMainInfo{}
	asinfo := model.AccountSecurityInfo{}
	atinfo := model.AccountThirdpartInfo{}
	alinfo := model.AccountLastloginInfo{}
	bFromLoginSer := strings.Contains(datareq.LoginRoute, "Login")
	if acResult, err := myredis.SelectAccount(accountdbName, datareq); err != nil {
		//如果是登录服加载
		if bFromLoginSer {
			row := dbrconn.Where("account = ? and passwd = ?", datareq.Account, datareq.Passwd).Find(&accountInfo)
			if row.RowsAffected < 1 {
				account.Errcode = msg.LoginErr_nouser
				return row.Error
			}
			if accountInfo.AccountStatus == 0 {
				account.Errcode = msg.LoginErr_accountforbidden
				return nil
			}
			datareq.Userid = accountInfo.UserID
		} else {
			row := dbrconn.Where("user_id = ? and passwd = ?", datareq.Userid, datareq.Passwd).Find(&accountInfo)
			if row.RowsAffected < 1 {
				account.Errcode = msg.LoginErr_nouser
				return row.Error
			}
			if accountInfo.AccountStatus == 0 {
				account.Errcode = msg.LoginErr_accountforbidden
				return nil
			}
		}
		dbrconn.Find(&asinfo, datareq.Userid)
		dbrconn.Find(&atinfo, datareq.Userid)
		dbrconn.Find(&alinfo, datareq.Userid)
		account.Token = base.GernateToken(account.Account, datareq.Userid)
		loginInfo.UserID = datareq.Userid
		//更新登录记录
		if bFromLoginSer {
			dbwconn.Save(&loginInfo)
		}
		atinfo.UserID = datareq.Userid
		asinfo.UserID = datareq.Userid
		myredis.WriteAccountInfo(accountdbName, account.Token, accountInfo, atinfo, asinfo, loginInfo)
		//写入redis
	} else {
		accountInfo = acResult[0].(model.AccountMainInfo)
		asinfo = acResult[1].(model.AccountSecurityInfo)
		atinfo = acResult[2].(model.AccountThirdpartInfo)
		alinfo = acResult[3].(model.AccountLastloginInfo)
		//重新生成token
		if len(acResult) == 5 {
			account.Token = acResult[4].(string)
			myredis.WriteAccountToken(accountdbName, account.Token, accountInfo.UserID)
		}
		updateAcountAllTTF([]string{
			fmt.Sprintf(myredis.ACCOUNT_MAIN_KEY, account.Account),
			fmt.Sprintf(myredis.ACCOUNT_THIRDINFO_KEY, account.Userid),
			fmt.Sprintf(myredis.ACCOUNT_SECURITY_KEY, account.Userid),
			fmt.Sprintf(myredis.ACCOUNT_LOGIN_KEY, account.Userid),
			fmt.Sprintf(myredis.ACCOUNT_TOKEN_KEY, account.Userid),
		}, account.Userid, []time.Duration{
			myredis.REDIS_ACCOUNT_TTF,
			myredis.REDIS_ACCOUNT_TTF,
			myredis.REDIS_ACCOUNT_TTF,
			myredis.REDIS_ACCOUNT_TTF,
			myredis.REDIS_TOKEN_TTF,
		})
		//if bFromLoginSer {

		//}

	}
	account.Userid = accountInfo.UserID
	account.AppID = accountInfo.AppID
	account.SiteID = accountInfo.SiteID
	account.Account = accountInfo.Account
	account.AcIndex = accountInfo.AccountIndex
	account.AcType = accountInfo.AccountType
	account.IsValid = accountInfo.AccountStatus
	account.NickName = accountInfo.NickName
	account.RegChan = accountInfo.ChannelID
	account.RegTime = accountInfo.RegisteDate
	account.RegSiteID = accountInfo.SiteID
	account.Gender = accountInfo.Gender
	account.ThirdHeadUrl = atinfo.HeadUrl //第三方头像地址
	account.Phonebinded = asinfo.PhoneNum //绑定手机号
	log.Debug("%v", account)
	return nil
}
func updateAcountAllTTF(key []string, userID int64, ttfs []time.Duration) {
	myredis.UpdateRedisTTF(accountdbName, userID, key, ttfs)
}
func UpdateATokenTTF(userID int64) {
	myredis.UpdateRedisTTF(accountdbName, userID, []string{fmt.Sprintf(myredis.ACCOUNT_TOKEN_KEY, userID)}, []time.Duration{myredis.REDIS_TOKEN_TTF})
}
