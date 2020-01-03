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

	// row := Accountconn.dbconn.Raw("select userid,Accounts, Nullity,TypeID,Gender, RegisterDate,ChannelID,SiteID from QPGameUserDB.dbo.AccountsInfo where UserID = ? ", datareq.Userid).Row()
	// err := row.Scan(&account.Userid, &account.Useraccount, &account.Accountfreezed, &account.Useraccounttype, &account.Usergender, &account.Userregistetime,
	// 	&account.Userregistechannelshort, &account.Userregistesiteid)
	// log.Debug("%v", err)
	// if err != nil {
	// 	account.Errcode = msg.LoginErr_nouser
	// 	return err
	// }
	// if account.Accountfreezed == true {
	// 	account.Errcode = msg.LoginErr_accountforbidden
	// 	return nil
	// }
	// row = Accountconn.dbconn.Raw("SELECT NickName,Mobile from QPGameUserDB.dbo.AccountsSecurityInfo where UserID = ? ", datareq.Userid).Row()
	// err = row.Scan(&account.Usernickname, &account.Userphonebinded)
	// row = Accountconn.dbconn.Raw("SELECT gamecoin,bankcoin from QPCoinDB.dbo.CoinData where UserID = ? ", datareq.Userid).Row()
	// err = row.Scan(&account.Usergamecoin, &account.Userbankcoin)
	// row = Accountconn.dbconn.Raw("SELECT RedPacket from QPGameUserDB.dbo.UserPropEx where UserID = ? ", datareq.Userid).Row()
	// err = row.Scan(&account.Usergoldbean)
	// row = Accountconn.dbconn.Raw("SELECT ExternalUserHeadImgUrl from QPGameExternalUserDB.dbo.ExternalUserMappingLocalRecord where LocalUserID = ? ", datareq.Userid).Row()
	// err = row.Scan(&account.Userthirdheadurl)
	// row = Accountconn.dbconn.Raw("SELECT VIPExperience from QPUserLobbyDB.dbo.AccountVIPProperty where LocalUserID = ? ", datareq.Userid).Row()
	// err = row.Scan(&account.VipExp)
	// row = Accountconn.dbconn.Raw("SELECT is_first_user from QPUserLobbyDB.dbo.User_Activity_MachineSerial where LocalUserID = ? ", datareq.Userid).Row()
	// err = row.Scan(&account.Useraccountindex)
	// row = Accountconn.dbconn.Raw("SELECT WIN_COUNT,LOST_COUNT,PLAY_TIME_COUNT,ONLINE_TIME_COUNT from QPGameTestDB.dbo.game_data where GAME_ID = ? and USER_ID = ? ", datareq.Gameid, datareq.Userid).Row()
	// err = row.Scan(&account.Gamewintimes, &account.Gamelosetimes, &account.Gameplaytime, &account.Gameonlinetime)
	// proplist := []msg.PropInfo{}
	// Accountconn.dbconn.Raw("SELECT propid,proptype,propnum, deadline FROM QPUserLobbyDB.dbo.UserLobbyProp where userid = ? and propid in (?)", datareq.Userid, datareq.Proplist).Scan(&proplist)
	// account.Proplist = append(account.Proplist, proplist...)
	// row = Accountconn.dbconn.Raw("SELECT gameid,serverid from QPGameUserDB.dbo.UserOnlineStatus where userid = ?", datareq.Userid).Row()
	// err = row.Scan(&account.Othergameid, &account.Otherroomid)
	// if account.Othergameid == 0 {
	// 	useronline := msg.Useronlinestatus{
	// 		Userid:      datareq.Userid,
	// 		Siteid:      account.Userregistesiteid,
	// 		Gameid:      datareq.Gameid,
	// 		Serverid:    datareq.Serverid,
	// 		Collectdate: time.Now(),
	// 	}
	// 	Accountconn.dbconn.Save(&useronline)
	// 	loginaddr := strings.Split(datareq.Loginip, ":")
	// 	Accountconn.dbconn.Exec("update QPGameUserDB.dbo.AccountsInfo set Lastlogondate = ?, Lastlogonip = ? where Userid = ? ", time.Now(), loginaddr[0], datareq.Userid)
	// 	//Raw方法只执行查询语句
	// 	Accountconn.dbconn.Exec("INSERT INTO QPGameUserDB.dbo.UserOnlineStatus (UserID,SiteID,GameID,ServerID) VALUES (?,?,?,?)", datareq.Userid, account.Userregistesiteid, datareq.Gameid, datareq.Serverid)
	// 	// log.Debug("%v", errs)
	// }
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
	log.Debug("%v", account)
	//tx.Commit()
	return nil
}
