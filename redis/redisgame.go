package myredis

import (
	"errors"

	"github.com/hudgit2019/leafboot/model"
	"github.com/hudgit2019/leafboot/msg"
)

const (
	GAMEDATA_KEY = "GameData:%d:%d"
)

func SelectGameData(redisName string, datareq *msg.Accountrdatareq) ([]interface{}, error) {
	gamedata := model.GameDataRecord{}
	return []interface{}{gamedata}, errors.New("SelectGameData no redis available")
}
func SaveGameData(redisName string, data msg.Accountwdatareq) error {
	//更新游戏属性
	return errors.New("SaveGameData no redis available")
}
