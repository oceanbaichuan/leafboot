package myredis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
	"github.com/hudgit2019/leafboot/base"
	"github.com/hudgit2019/leafboot/model"
	"github.com/hudgit2019/leafboot/msg"
	"github.com/name5566/leaf/log"
)

const (
	ACCOUNT_ACCOUNTMAIN_KEY = "AccountMain:%s"
	ACCOUNT_USERIDMAIN_KEY  = "UserIDMain:%d" //关联账号名
	ACCOUNT_LOGIN_KEY       = "AccountLogin:%d"
	ACCOUNT_TOKEN_KEY       = "AccountToken:%d"
	ACCOUNT_THIRDINFO_KEY   = "AccountThirdInfo:%d"
	ACCOUNT_SECURITY_KEY    = "AccountSecurity:%d"
)

func SelectAccount(redisName string, datareq *msg.Accountrdatareq) ([]interface{}, error) {
	rdconn, err := GetRedis(redisName, datareq.Userid)
	if err != nil {
		log.Error("redis SelectAccount GetRedis err:%v", err)
		return nil, err
	}
	//判断token是否有效
	//bLoginServer := strings.Contains(datareq.LoginRoute, "Login")
	bHasToken := false
	//首先检验token
	if datareq.Userid > 0 {
		accvalue, err := rdconn.Get(fmt.Sprintf(ACCOUNT_TOKEN_KEY, datareq.Userid)).Result()
		if err == nil {
			if accvalue != datareq.Token {
				strErr := fmt.Sprintf("redis SelectAccount userID:%d clienttoken:%s redistoken:%s",
					datareq.Userid, datareq.Token, accvalue)
				log.Error(strErr)
				return nil, errors.New("token is wrong")
			}
			bHasToken = true
		} else {
			log.Error("redis SelectAccount err:%v", err)
		}
		accvalue, err = rdconn.Get(fmt.Sprintf(ACCOUNT_USERIDMAIN_KEY, datareq.Userid)).Result()
		if err != nil {
			return nil, errors.New("user is not found")
		}
		//uid的account覆盖客户端给的，避免用户串号
		datareq.Account = accvalue
	}
	accvalue, err := rdconn.Get(fmt.Sprintf(ACCOUNT_ACCOUNTMAIN_KEY, datareq.Account)).Result()
	if err != nil {
		log.Debug("redis SelectAccount account:%s err:%v", datareq.Account, err)
		return nil, err
	}
	accountInfo := model.AccountMainInfo{}
	asinfo := model.AccountSecurityInfo{}
	atinfo := model.AccountThirdpartInfo{}
	alinfo := model.AccountLastloginInfo{}
	token := ""
	err = json.Unmarshal([]byte(accvalue), &accountInfo)
	if err != nil {
		return nil, err
	}
	//token匹配失败时，需要匹配密码
	if !bHasToken && accountInfo.Passwd != datareq.Passwd {
		strErr := fmt.Sprintf("redis SelectAccount userID:%d clientpasswd:%s redispasswd:%s",
			datareq.Userid, datareq.Passwd, accountInfo.Passwd)
		log.Error(strErr)
		return nil, errors.New("passwd is wrong")
	} else if !bHasToken {
		token = base.GernateToken(datareq.Account, accountInfo.UserID)
	}
	strThirdKey := fmt.Sprintf(ACCOUNT_THIRDINFO_KEY, accountInfo.UserID)
	strSecKey := fmt.Sprintf(ACCOUNT_SECURITY_KEY, accountInfo.UserID)
	strLoginKey := fmt.Sprintf(ACCOUNT_LOGIN_KEY, accountInfo.UserID)
	cmds, err := rdconn.Pipelined(func(accPipe redis.Pipeliner) error {
		//读取第三方信息
		accPipe.Get(strThirdKey)
		//读取个人隐私信息
		accPipe.Get(strSecKey)
		accPipe.Get(strLoginKey)
		return nil
	})
	if err != nil {
		log.Debug("redis SelectAccount Pipelined err:%v", err)
		return nil, err
	}
	for _, v := range cmds {
		if v.Err() != nil {
			log.Error("redis SelectAccount cmds key:%s err:%v", v.Args(), err)
			return nil, v.Err()
		}
		if strings.Contains(v.Args()[1].(string), strThirdKey) {
			//读取第三方信息
			accvalue, err = v.(*redis.StringCmd).Result()
			if err != nil {
				log.Error("redis SelectAccount err:%v", err)
				return nil, err
			} else {
				err = json.Unmarshal([]byte(accvalue), &atinfo)
				if err == nil {
				} else {
					log.Error("redis SelectAccount Unmarshal atinfo  err:%v", err)
					return nil, err
				}
			}
		} else if strings.Contains(v.Args()[1].(string), strSecKey) {
			//读取个人隐私信息
			accvalue, err = v.(*redis.StringCmd).Result()
			if err != nil {
				log.Error("redis SelectAccount err:%v", err)
				return nil, err
			} else {
				err = json.Unmarshal([]byte(accvalue), &asinfo)
				if err == nil {
				} else {
					log.Error("redis SelectAccount Unmarshal asinfo err:%v", err)
					return nil, err
				}
			}
		} else if strings.Contains(v.Args()[1].(string), strLoginKey) {
			//读取上次登录信息
			accvalue, err = v.(*redis.StringCmd).Result()
			if err != nil {
				log.Error("redis SelectAccount err:%v", err)
				return nil, err
			} else {
				err = json.Unmarshal([]byte(accvalue), &alinfo)
				if err == nil {
				} else {
					log.Error("redis SelectAccount Unmarshal alinfo err:%v", err)
					return nil, err
				}
			}
		}
	}
	return []interface{}{accountInfo, asinfo, atinfo, alinfo, token}, nil
}
func WriteAccountInfo(redisName string, token string, aInfo model.AccountMainInfo,
	aThirdInfo model.AccountThirdpartInfo, aSecurityInfo model.AccountSecurityInfo, aLoginInfo model.AccountLastloginInfo) {
	rdconn, err := GetRedis(redisName, aInfo.UserID)
	if err != nil {
		log.Error("redis SelectUserInfo err:%v", err)
		return
	}
	rdconn.Pipelined(func(accPipe redis.Pipeliner) error {
		propStr, err := json.Marshal(&aInfo)
		if err == nil {
			accPipe.SetNX(fmt.Sprintf(ACCOUNT_ACCOUNTMAIN_KEY, aInfo.Account), propStr, REDIS_ACCOUNT_TTF)
		}
		propStr, err = json.Marshal(&aThirdInfo)
		if err == nil {
			accPipe.SetNX(fmt.Sprintf(ACCOUNT_THIRDINFO_KEY, aThirdInfo.UserID), propStr, REDIS_ACCOUNT_TTF)
		}
		propStr, err = json.Marshal(&aSecurityInfo)
		if err == nil {
			accPipe.SetNX(fmt.Sprintf(ACCOUNT_SECURITY_KEY, aSecurityInfo.UserID), propStr, REDIS_ACCOUNT_TTF)
		}
		propStr, err = json.Marshal(&aLoginInfo)
		if err == nil {
			accPipe.SetNX(fmt.Sprintf(ACCOUNT_LOGIN_KEY, aLoginInfo.UserID), propStr, REDIS_ACCOUNT_TTF)
		}
		accPipe.SetNX(fmt.Sprintf(ACCOUNT_TOKEN_KEY, aLoginInfo.UserID), token, REDIS_TOKEN_TTF)
		return nil
	})
}

func WriteAccountToken(redisName string, token string, userID int64) {
	rdconn, err := GetRedis(redisName, userID)
	if err != nil {
		log.Error("redis SelectUserInfo err:%v", err)
		return
	}
	rdconn.Pipelined(func(accPipe redis.Pipeliner) error {
		accPipe.SetNX(fmt.Sprintf(ACCOUNT_TOKEN_KEY, userID), token, REDIS_TOKEN_TTF)
		return nil
	})
}
