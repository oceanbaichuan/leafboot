package myredis

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/hudgit2019/leafboot/conf"
	"github.com/name5566/leaf/log"
)

type RedisShareding struct {
	dbinfo conf.RedisInfo
	dbconn *redis.Client
}

const (
	REDIS_TOKEN_TTF   = 60 * 30 * time.Second       //token有效期30min
	REDIS_ACCOUNT_TTF = 3600 * 24 * 7 * time.Second //账号属性有效期7day
)

//mapDBList 数据库池。key为数据库实例用途名，例如:AccountDB,GameUserDB,GameDataDB,LogDB
var mapDBList map[string][]RedisShareding

func init() {
	mapDBList = make(map[string][]RedisShareding)

}
func StartRedis() {
	for _, v := range conf.Server.RedisList {
		err := OpenRedisGroup(v)
		if err != nil {
			log.Fatal(" StartRedis err:%v", err)

		}
	}
}
func OpenRedisGroup(dbinfo conf.RedisInfo) error {
	if _, ok := conf.Server.CustomDBName[dbinfo.RedisName]; !ok {
		log.Error("OpenRedisGroup RedisName:%s not supported", dbinfo.RedisName)
		return nil
	}
	bHas := false
	for k, rslist := range mapDBList {
		for i, v := range rslist {
			if v.dbinfo.Addr == dbinfo.Addr &&
				v.dbinfo.RedisName == dbinfo.RedisName {
				mapDBList[k][i].dbinfo = dbinfo
				bHas = true
				break
			}
		}
		if bHas {
			break
		}
	}
	if bHas {
		return nil
	}
	tmpconn, err := OpenRedis(dbinfo)
	if err == nil {
		dbSharding := RedisShareding{
			dbinfo: dbinfo,
			dbconn: tmpconn,
		}
		if strings.Contains(dbinfo.RedisName, "plat_account_db") {
			mapDBList["plat_account_db"] = append(mapDBList["plat_account_db"], dbSharding)
		} else if strings.Contains(dbinfo.RedisName, "game_user_db") {
			mapDBList["game_user_db"] = append(mapDBList["game_user_db"], dbSharding)
		} else if strings.Contains(dbinfo.RedisName, "game_data_db") {
			mapDBList["game_data_db"] = append(mapDBList["game_data_db"], dbSharding)
		} else if strings.Contains(dbinfo.RedisName, "game_log_db") {
			mapDBList["game_log_db"] = append(mapDBList["game_log_db"], dbSharding)
		}
	} else {
		return err
	}
	return nil
}

func OpenRedis(dbinfo conf.RedisInfo) (dbconn *redis.Client, err error) {
	redisdb := redis.NewClient(&redis.Options{
		Addr:     dbinfo.Addr,
		Password: dbinfo.Passwd,
		DB:       int(dbinfo.Slot),
	})
	_, err = redisdb.Ping().Result()
	if err != nil {
		return nil, err
	}
	return redisdb, nil
}
func GetRedis(dbName string, userID int64) (*redis.Client, error) {
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
func UpdateRedisTTF(redisName string, userID int64, key []string, ttfs []time.Duration) {
	rdconn, err := GetRedis(redisName, userID)
	if err != nil {
		log.Error("redis UpdateATTF err:%v", err)
		return
	}
	rdconn.Pipelined(func(accPipe redis.Pipeliner) error {
		for i, v := range key {
			accPipe.Expire(v, ttfs[i])
		}
		return nil
	})
	//rdconn.Expire(key, ttf)
}
