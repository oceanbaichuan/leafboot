package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hudgit2019/leafboot/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/name5566/leaf/log"
)

type DBShareding struct {
	dbinfo conf.DatabaseInfo
	dbconn *gorm.DB
}

//mapDBList 数据库池。key为数据库实例用途名，例如:AccountDB,GameUserDB,GameDataDB,LogDB
var mapDBList map[string][]DBShareding

func init() {
	mapDBList = make(map[string][]DBShareding)

}
func StartDB() {
	for _, v := range conf.Server.DbList {
		tmpconn, err := OpenDB(v)
		if err == nil {
			dbSharding := DBShareding{
				dbinfo: v,
				dbconn: tmpconn,
			}
			if strings.Contains(v.DataBase, "plat_account_db") {
				mapDBList["plat_account_db"] = append(mapDBList["plat_account_db"], dbSharding)
			} else if strings.Contains(v.DataBase, "game_user_db") {
				mapDBList["game_user_db"] = append(mapDBList["game_user_db"], dbSharding)
			} else if strings.Contains(v.DataBase, "game_data_db") {
				mapDBList["game_data_db"] = append(mapDBList["game_data_db"], dbSharding)
			} else if strings.Contains(v.DataBase, "game_log_db") {
				mapDBList["game_log_db"] = append(mapDBList["game_log_db"], dbSharding)
			}
		} else {
			log.Fatal(" StartDB err:%v", err)
		}
	}
}
func OpenDB(dbinfo conf.DatabaseInfo) (dbconn *gorm.DB, err error) {

	if dbinfo.DbType == "mysql" {
		strActDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			dbinfo.UserName, dbinfo.Passwd, dbinfo.Host, dbinfo.Port, dbinfo.DataBase)
		dbconn, err := gorm.Open("mysql", strActDB)
		if err == nil {
			dbconn.SingularTable(true)
			dbconn.DB().SetMaxIdleConns(10)
			dbconn.DB().SetMaxOpenConns(100)
		}
		return dbconn, err
	} else if dbinfo.DbType == "mssql" {
		strActDB := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=1433;encrypt=disable;parseTime=True",
			dbinfo.Host, dbinfo.DataBase, dbinfo.UserName, dbinfo.Passwd)
		dbconn, err := gorm.Open("mssql", strActDB)
		if err == nil {
			dbconn.SingularTable(true)
			dbconn.DB().SetMaxIdleConns(10)
			dbconn.DB().SetMaxOpenConns(100)
		}
		return dbconn, err
	}
	return nil, errors.New("dbType not supported")
}
func GetDB(dbName string, userID int64) (*gorm.DB, error) {
	if v, ok := mapDBList[dbName]; ok {
		for _, dbSharding := range v {
			if dbSharding.dbinfo.MinUID <= userID && dbSharding.dbinfo.MaxUID >= userID { //在分区内
				return dbSharding.dbconn, nil
			} else if dbSharding.dbinfo.MinUID == dbSharding.dbinfo.MaxUID { //不限分区
				return dbSharding.dbconn, nil
			} else if dbSharding.dbinfo.MinUID <= userID && dbSharding.dbinfo.MaxUID == -1 { //上不封顶
				return dbSharding.dbconn, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("dbName:%s not found", dbName))
}
