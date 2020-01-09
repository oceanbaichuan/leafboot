package myredis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/hudgit2019/leafboot/model"
	"github.com/hudgit2019/leafboot/msg"
	"github.com/name5566/leaf/log"
)

const (
	USER_CURRENCY_KEY     = "UserCurrency:%d"       //UserCurrency:userid 货币
	USER_BACKPACK_SET_KEY = "UserBackPack:%d"       //UserBackPack:userid:propid 所有道具ID集合
	USER_PROPID_SET_KEY   = "UserBackPack:%v:%v"    //UserBackPack:userid:propid 道具ID的type集合
	USER_PROPINFO_KEY     = "UserBackPack:%v:%v:%v" //UserBackPack:userid:propid:proptype 道具详细
	USER_PROPERTY_KEY     = "UserProperty:%d"       //UserProperty:userid 用户平台属性
	USER_DAYPROPERTY_KEY  = "UserDayProperty:%d"    //UserDayProperty:userid 用户每日平台属性
	USER_ONLINE_KEY       = "UserOnline:%d"         //UserOnline:userid 用户在线
)

func SelectUserInfo(redisName string, datareq *msg.Accountrdatareq) ([]interface{}, error) {
	rdconn, err := GetRedis(redisName, datareq.Userid)
	if err != nil {
		log.Error("redis SelectUserInfo err:%v", err)
		return nil, err
	}
	userCurrency := model.GameUserCurrency{}
	userProperty := model.GameUserProperty{}
	userDayProperty := model.GameUserDayProperty{}
	userOnline := model.GameUserOnline{}
	var tmpPropList []msg.PropInfo
	var propIDList []string
	propTypeList := make(map[string][]string)
	strCurrency := fmt.Sprintf(USER_CURRENCY_KEY, datareq.Userid)
	strProperty := fmt.Sprintf(USER_PROPERTY_KEY, datareq.Userid)
	strDayProperty := fmt.Sprintf(USER_DAYPROPERTY_KEY, datareq.Userid)
	strOnline := fmt.Sprintf(USER_ONLINE_KEY, datareq.Userid)
	strBackPakc := "UserBackPack"
	if len(datareq.Proplist) > 0 {
		for _, propid := range datareq.Proplist {
			propTypeList[strconv.Itoa(int(propid))], err = rdconn.ZRange(fmt.Sprintf(USER_PROPID_SET_KEY, datareq.Userid, propid), 0, -1).Result()
			if err != nil {
				log.Error("redis SelectUserInfo propTypeList err:%v", err)
			}
		}
	} else {
		//加载所有道具
		propIDList, err = rdconn.ZRange(fmt.Sprintf(USER_BACKPACK_SET_KEY, datareq.Userid), 0, -1).Result()
		if err != nil {
			log.Error("redis SelectUserInfo propIDList err:%v", err)
		} else {
			for _, propid := range propIDList {
				propTypeList[propid], err = rdconn.ZRange(fmt.Sprintf(USER_PROPID_SET_KEY, datareq.Userid, propid), 0, -1).Result()
			}
		}
	}
	onlineStr, err := rdconn.Get(strOnline).Result()
	json.Unmarshal([]byte(onlineStr), &userOnline)
	cmds, err := rdconn.Pipelined(func(accPipe redis.Pipeliner) error {
		accPipe.HGetAll(strCurrency)
		//accPipe.ObjectEncoding(strCurrency)
		accPipe.HGetAll(strProperty)
		accPipe.HGetAll(strDayProperty)
		log.Debug("redisuser datareq.Proplist:%v propTypeList:%v", datareq.Proplist, propTypeList)
		for k, typelist := range propTypeList {
			for _, v := range typelist {
				accPipe.HGetAll(fmt.Sprintf(USER_PROPINFO_KEY, datareq.Userid, k, v))
			}
		}
		//accPipe.Get(strOnline)
		return nil
	})
	if err != nil {
		log.Error("redis Pipelined err:%v", err)
		return nil, err
	}
	for _, v := range cmds {
		if v.Err() != nil {
			log.Error("redis cmds name:%s err:%v", v.Name(), v.Err())
			return nil, v.Err()
		}
		//log.Debug("redis cmds v:%v v.Name:%v v.Args(1):%v", v, v.Name(), v.Args()[1])
		if strings.Contains(v.Args()[1].(string), strCurrency) {
			mapCmd, _ := v.(*redis.StringStringMapCmd).Result()
			if len(mapCmd) == 0 {
				log.Error("redis strProperty err:no data")
				return nil, errors.New("Currency no data")
			}
			userCurrency.UpdateTime, _ = time.Parse("01/02/2006", mapCmd["update_time"])
			userCurrency.GameCoin, _ = strconv.ParseInt(mapCmd["game_coin"], 10, 64)
			tmpPrize, _ := strconv.Atoi(mapCmd["prize_ticket"])
			userCurrency.PrizeTicket = int32(tmpPrize)
		} else if strings.Contains(v.Args()[1].(string), strProperty) {
			mapCmd, err := v.(*redis.StringStringMapCmd).Result()
			if err != nil {
				log.Error("redis strProperty err:%v", err)
				return nil, err
			}
			if len(mapCmd) == 0 {
				log.Error("redis strProperty err:no data")
				return nil, errors.New("Property no data")
			}
			userProperty.UpdateTime, _ = time.Parse("01/02/2006", mapCmd["update_time"])
			userProperty.PlayCoin, _ = strconv.ParseInt(mapCmd["play_coin"], 10, 64)
			userProperty.PlayTax, _ = strconv.ParseUint(mapCmd["play_tax"], 10, 64)
			tmpPrize, _ := strconv.Atoi(mapCmd["play_prizeticket"])
			userProperty.PlayPrizeticket = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["play_time"])
			userProperty.PlayTime = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["online_time"])
			userProperty.OnlineTime = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["help_times"])
			tmpPrize, _ = strconv.Atoi(mapCmd["recharge_times"])
			userProperty.RechargeTimes = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["recharge_money"])
			userProperty.RechargeMoney = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["vip_exp"])
			userProperty.VipExp = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["game_exp"])
			userProperty.GameExp = int32(tmpPrize)
		} else if strings.Contains(v.Args()[1].(string), strDayProperty) {
			mapCmd, err := v.(*redis.StringStringMapCmd).Result()
			//log.Debug("redis mapCmd:%v err:%v", mapCmd, err)
			if err != nil {
				log.Error("redis strDayProperty err:%v", err)
				return nil, err
			}
			if len(mapCmd) == 0 {
				log.Error("redis strDayProperty err:no data")
				return nil, errors.New("DayProperty no data")
			}
			userDayProperty.UpdateTime, _ = time.Parse("01/02/2006", mapCmd["update_time"])
			userDayProperty.PlayCoin, _ = strconv.ParseInt(mapCmd["play_coin"], 10, 64)
			userDayProperty.PlayTax, _ = strconv.ParseUint(mapCmd["play_tax"], 10, 64)
			tmpPrize, _ := strconv.Atoi(mapCmd["play_prizeticket"])
			userDayProperty.PlayPrizeticket = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["play_time"])
			userDayProperty.PlayTime = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["online_time"])
			userDayProperty.OnlineTime = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["recharge_times"])
			userDayProperty.RechargeTimes = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["recharge_money"])
			userDayProperty.RechargeMoney = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["help_times"])
			userDayProperty.HelpTimes = int32(tmpPrize)
		} else if strings.Contains(v.Args()[1].(string), strBackPakc) {
			mapCmd, _ := v.(*redis.StringStringMapCmd).Result()
			//log.Debug("redis mapCmd:%v err:%v", mapCmd, err)
			// if err != nil {
			// 	log.Error("redis strBackPakc err:%v", err)
			// 	return nil, err
			// }
			propInfo := msg.PropInfo{}
			propInfo.Proptime, _ = strconv.ParseInt(mapCmd["prop_time"], 10, 64)
			tmpPrize, _ := strconv.Atoi(mapCmd["prop_id"])
			propInfo.Propid = int32(tmpPrize)
			tmpPrize, _ = strconv.Atoi(mapCmd["prop_type"])
			propInfo.Proptype = int32(tmpPrize)
			tmpPropList = append(tmpPropList, propInfo)
		}
	}
	return []interface{}{userCurrency, userProperty, userDayProperty, tmpPropList, userOnline}, nil
}

//WriteUserInfoRedis 首次登陆，信息写入redis
func WriteUserInfoRedis(redisName string, userProperty model.GameUserProperty, userDayProperty model.GameUserDayProperty,
	userCurrency model.GameUserCurrency, propList []model.GameUserBackpack, userOnline model.GameUserOnline) {

	rdconn, err := GetRedis(redisName, userProperty.UserID)
	if err != nil {
		log.Error("redis SelectUserInfo err:%v", err)
		return
	}
	propStr, _ := json.Marshal(&userProperty)
	var mapProp map[string]interface{}
	json.Unmarshal([]byte(propStr), &mapProp)
	dayPropStr, _ := json.Marshal(&userDayProperty)
	var mapDayProp map[string]interface{}
	json.Unmarshal([]byte(dayPropStr), &mapDayProp)
	currencyStr, _ := json.Marshal(&userCurrency)
	var mapCurrency map[string]interface{}
	json.Unmarshal([]byte(currencyStr), &mapCurrency)
	onlineStr, _ := json.Marshal(&userOnline)
	var mapOnline map[string]interface{}
	json.Unmarshal([]byte(onlineStr), &mapOnline)
	var backPackStr []map[string]interface{}
	for _, v := range propList {
		propStr, err := json.Marshal(&v)
		if err == nil {
			var mapProp map[string]interface{}
			json.Unmarshal([]byte(propStr), &mapProp)
			backPackStr = append(backPackStr, mapProp)
		}
	}
	rdconn.Pipelined(func(accPipe redis.Pipeliner) error {
		accPipe.HMSet(fmt.Sprintf(USER_CURRENCY_KEY, userProperty.UserID), mapCurrency)
		accPipe.Expire(fmt.Sprintf(USER_CURRENCY_KEY, userProperty.UserID), REDIS_ACCOUNT_TTF)
		accPipe.HMSet(fmt.Sprintf(USER_PROPERTY_KEY, userProperty.UserID), mapProp)
		accPipe.Expire(fmt.Sprintf(USER_PROPERTY_KEY, userProperty.UserID), REDIS_ACCOUNT_TTF)
		accPipe.HMSet(fmt.Sprintf(USER_DAYPROPERTY_KEY, userProperty.UserID), mapDayProp)
		accPipe.Expire(fmt.Sprintf(USER_DAYPROPERTY_KEY, userProperty.UserID), REDIS_ACCOUNT_TTF)
		accPipe.HMSet(fmt.Sprintf(USER_ONLINE_KEY, userProperty.UserID), mapOnline)
		for _, v := range backPackStr {
			log.Debug("WriteUserInfoRedis backPack:%v", v)
			accPipe.HMSet(fmt.Sprintf(USER_PROPINFO_KEY, userProperty.UserID, v["prop_id"], v["prop_type"]), v)
			accPipe.Expire(fmt.Sprintf(USER_PROPINFO_KEY, userProperty.UserID, v["prop_id"], v["prop_type"]), REDIS_ACCOUNT_TTF)
			accPipe.ZAdd(fmt.Sprintf(USER_BACKPACK_SET_KEY, userProperty.UserID),
				redis.Z{Score: 100, Member: v["prop_id"]})
			accPipe.Expire(fmt.Sprintf(USER_BACKPACK_SET_KEY, userProperty.UserID), REDIS_ACCOUNT_TTF)
			accPipe.ZAdd(fmt.Sprintf(USER_PROPID_SET_KEY, userProperty.UserID, v["prop_id"]),
				redis.Z{Score: 100, Member: v["prop_type"]})
			accPipe.Expire(fmt.Sprintf(USER_PROPID_SET_KEY, userProperty.UserID, v["prop_id"]), REDIS_ACCOUNT_TTF)
		}
		return nil
	})
}

func WriteUserProperty(redisName string, userProperty model.GameUserProperty) {
	rdconn, err := GetRedis(redisName, userProperty.UserID)
	if err != nil {
		log.Error("redis SelectUserInfo err:%v", err)
		return
	}
	propStr, err := json.Marshal(&userProperty)
	if err == nil {
		var mapProp map[string]interface{}
		json.Unmarshal([]byte(propStr), &mapProp)
		rdconn.HMSet(fmt.Sprintf(USER_PROPERTY_KEY, userProperty.UserID), mapProp)
	}

}
func WriteUserDayProperty(redisName string, userProperty model.GameUserDayProperty) {
	rdconn, err := GetRedis(redisName, userProperty.UserID)
	if err != nil {
		log.Error("redis SelectUserInfo err:%v", err)
		return
	}
	propStr, err := json.Marshal(&userProperty)
	if err == nil {
		var mapProp map[string]interface{}
		json.Unmarshal([]byte(propStr), &mapProp)
		rdconn.HMSet(fmt.Sprintf(USER_DAYPROPERTY_KEY, userProperty.UserID), mapProp)
	}
}
func WriteUserCurrency(redisName string, userProperty model.GameUserCurrency) {
	rdconn, err := GetRedis(redisName, userProperty.UserID)
	if err != nil {
		log.Error("redis SelectUserInfo err:%v", err)
		return
	}
	propStr, err := json.Marshal(&userProperty)
	if err == nil {
		var mapProp map[string]interface{}
		json.Unmarshal([]byte(propStr), &mapProp)
		rdconn.HMSet(fmt.Sprintf(USER_CURRENCY_KEY, userProperty.UserID), mapProp)
	}
}
func WriteUserOnline(redisName string, userProperty model.GameUserOnline) {
	rdconn, err := GetRedis(redisName, userProperty.UserID)
	if err != nil {
		log.Error("redis SelectUserInfo err:%v", err)
		return
	}
	propStr, err := json.Marshal(&userProperty)
	if err == nil {
		var mapProp map[string]interface{}
		json.Unmarshal([]byte(propStr), &mapProp)
		rdconn.HMSet(fmt.Sprintf(USER_ONLINE_KEY, userProperty.UserID), mapProp)
	}
}

func WritePropList(redisName string, userID int64, propList []model.GameUserBackpack) {
	rdconn, err := GetRedis(redisName, userID)
	if err != nil {
		log.Error("redis SelectUserInfo err:%v", err)
		return
	}
	var backPackStr []map[string]interface{}
	for _, v := range propList {
		propStr, err := json.Marshal(&v)
		if err == nil {
			var mapProp map[string]interface{}
			json.Unmarshal([]byte(propStr), &mapProp)
			backPackStr = append(backPackStr, mapProp)
		}
	}
	rdconn.Pipelined(func(accPipe redis.Pipeliner) error {
		for _, v := range backPackStr {
			accPipe.HMSet(fmt.Sprintf(USER_PROPINFO_KEY, userID, v["prop_id"], v["prop_type"]), v)
			accPipe.ZAdd(fmt.Sprintf(USER_BACKPACK_SET_KEY, userID),
				redis.Z{Score: v["prop_id"].(float64), Member: v["prop_id"]})
			accPipe.ZAdd(fmt.Sprintf(USER_PROPID_SET_KEY, userID, v["prop_id"]),
				redis.Z{Score: v["prop_type"].(float64), Member: v["prop_type"]})
		}
		return nil
	})
}
func SaveUserProperty(redisName string, data msg.Accountwdatareq) error {
	rdconn, err := GetRedis(redisName, data.UserID)
	if err != nil {
		log.Error("redis SelectUserInfo err:%v", err)
		return err
	}
	_, err = rdconn.Pipelined(func(accPipe redis.Pipeliner) error {
		//更新货币,默认redis都存在数据
		accPipe.HIncrBy(fmt.Sprintf(USER_CURRENCY_KEY, data.UserID), "game_coin", data.GameCoin)
		accPipe.HIncrBy(fmt.Sprintf(USER_CURRENCY_KEY, data.UserID), "prize_ticket", int64(data.PrizeTicket))
		//更新道具
		for _, v := range data.PropList {
			if v.Propnum != 0 {
				accPipe.HIncrBy(fmt.Sprintf(USER_PROPINFO_KEY, data.UserID, v.Propid, v.Proptype), "prop_count", int64(v.Propnum))
			}
			if v.Proptime != 0 {
				accPipe.HIncrBy(fmt.Sprintf(USER_BACKPACK_SET_KEY, data.UserID), "prop_time", v.Proptime)
			}
			accPipe.ZAdd(fmt.Sprintf(USER_BACKPACK_SET_KEY, data.UserID),
				redis.Z{Score: float64(v.Propid), Member: v.Propid})
			accPipe.ZAdd(fmt.Sprintf(USER_PROPID_SET_KEY, data.UserID, v.Propid),
				redis.Z{Score: float64(v.Proptype), Member: v.Proptype})
		}
		//更新平台属性
		// online_time
		// play_time
		// play_coin
		// play_tax
		// play_prizeticket
		// recharge_times
		// recharge_money
		// vip_exp
		// game_exp

		accPipe.HIncrBy(fmt.Sprintf(USER_PROPERTY_KEY, data.UserID), "play_time", int64(data.GamePlaytime))
		accPipe.HIncrBy(fmt.Sprintf(USER_PROPERTY_KEY, data.UserID), "online_time", int64(data.GameOnlinetime))
		accPipe.HIncrBy(fmt.Sprintf(USER_PROPERTY_KEY, data.UserID), "play_coin", int64(data.GameCoin))
		accPipe.HIncrBy(fmt.Sprintf(USER_PROPERTY_KEY, data.UserID), "play_tax", int64(data.GameTax))
		accPipe.HIncrBy(fmt.Sprintf(USER_PROPERTY_KEY, data.UserID), "game_exp", int64(data.GameExp))
		accPipe.HIncrBy(fmt.Sprintf(USER_PROPERTY_KEY, data.UserID), "play_prizeticket", int64(data.PrizeTicket))
		//更新平台日属性

		accPipe.HIncrBy(fmt.Sprintf(USER_DAYPROPERTY_KEY, data.UserID), "play_time", int64(data.GamePlaytime))
		accPipe.HIncrBy(fmt.Sprintf(USER_DAYPROPERTY_KEY, data.UserID), "online_time", int64(data.GameOnlinetime))
		accPipe.HIncrBy(fmt.Sprintf(USER_DAYPROPERTY_KEY, data.UserID), "play_coin", int64(data.GameCoin))
		accPipe.HIncrBy(fmt.Sprintf(USER_DAYPROPERTY_KEY, data.UserID), "play_tax", int64(data.GameTax))
		accPipe.HIncrBy(fmt.Sprintf(USER_DAYPROPERTY_KEY, data.UserID), "game_exp", int64(data.GameExp))
		accPipe.HIncrBy(fmt.Sprintf(USER_DAYPROPERTY_KEY, data.UserID), "play_prizeticket", int64(data.PrizeTicket))
		accPipe.HIncrBy(fmt.Sprintf(USER_DAYPROPERTY_KEY, data.UserID), "help_times", int64(data.HelpTimes))
		return nil
	})
	return err
}
func DeleteOnline(redisName string, userid int64) {
	rdconn, err := GetRedis(redisName, userid)
	if err != nil {
		log.Error("redis SelectUserInfo err:%v", err)
		return
	}
	rdconn.Del(fmt.Sprintf(USER_ONLINE_KEY, userid))
}
