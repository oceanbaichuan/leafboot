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

const (
	DBFLAG_RW = "write"    //可读写
	DBFLAG_R  = "readonly" //只读
)

//mapDBList 数据库池。key为数据库实例用途名，例如:AccountDB,GameUserDB,GameDataDB,LogDB
var mapMasterRDBList map[string][]DBShareding
var mapMasterRWDBList map[string][]DBShareding
var mapSlaveRDBList map[string][]DBShareding

func init() {
	mapMasterRDBList = make(map[string][]DBShareding)
	mapMasterRWDBList = make(map[string][]DBShareding)
	mapSlaveRDBList = make(map[string][]DBShareding)
}
func StartDB() error {
	for _, v := range conf.Server.DbList {
		err := OpenDBGroup(v)
		if err != nil {
			log.Error(" StartDB err:%v", err)
			return err
		}
	}
	return nil
}
func OpenDBGroup(dbinfo conf.DatabaseInfo) error {
	//M读写库
	if dbinfo.IsMaster == 1 &&
		strings.Contains(dbinfo.DbRWFlag, DBFLAG_RW) {
		return checkDB(mapMasterRWDBList, dbinfo)
	} else if dbinfo.IsMaster == 1 &&
		strings.Contains(dbinfo.DbRWFlag, DBFLAG_R) {
		//M只读库
		return checkDB(mapMasterRDBList, dbinfo)
	} else if dbinfo.IsMaster == 0 {
		//S库
		return checkDB(mapSlaveRDBList, dbinfo)
	}
	return nil
}

//添加DB
func addDB(mapDb map[string][]DBShareding, dbName string, sharding DBShareding) error {
	mapDb[dbName] = append(mapDb[dbName], sharding)
	return nil
}

//checkDB检验db是否存在，存在则更新，不存在则创建
func checkDB(mapDb map[string][]DBShareding, dbinfo conf.DatabaseInfo) error {
	bHas := false
	for k, rslist := range mapDb {
		for i, v := range rslist {
			if v.dbinfo.Host == dbinfo.Host &&
				v.dbinfo.Port == dbinfo.Port &&
				v.dbinfo.DataBase == dbinfo.DataBase {
				mapDb[k][i].dbinfo = dbinfo
				bHas = true
				break
			}
		}
		if bHas {
			break
		}
	}
	if !bHas {
		tmpconn, err := OpenDB(dbinfo)
		if err == nil {
			dbSharding := DBShareding{
				dbinfo: dbinfo,
				dbconn: tmpconn,
			}
			return addDB(mapDb, dbinfo.DataBase, dbSharding)
		} else {
			return err
		}
	}
	return nil
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

//GetDB dbName:数据库名 rwFlag:write,readonly
func GetDB(dbName string, userID int64, rwFlag string, bMaster bool) (*gorm.DB, error) {
	var dbconn *gorm.DB = nil
	var err error = nil
	if bMaster && strings.Contains(rwFlag, "write") {
		dbconn, err = getMasterRW(dbName, userID)

	} else if bMaster && strings.Contains(rwFlag, "readonly") {
		dbconn, err = getMasterR(dbName, userID)
	}
	if dbconn == nil {
		dbconn, err = getSlaveR(dbName, userID)
	}
	return dbconn, err
}
func getMasterRW(dbName string, userID int64) (*gorm.DB, error) {
	if v, ok := mapMasterRWDBList[dbName]; ok {
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
func getMasterR(dbName string, userID int64) (*gorm.DB, error) {
	if v, ok := mapMasterRDBList[dbName]; ok {
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
func getSlaveR(dbName string, userID int64) (*gorm.DB, error) {
	if v, ok := mapSlaveRDBList[dbName]; ok {
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
