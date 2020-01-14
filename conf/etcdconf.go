package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
	"go.etcd.io/etcd/client"
)

//全平台统一，应用层无需修改
const (
	GrayPathPrefix   = "/GrayConf/"
	DBEtcdDir        = "/DBConf/"
	RedisEtcdDir     = "/RedisConf/"
	FrontendDir      = "/Frontend/"   //前台目录
	MiddlePlatDir    = "/MiddlePlat/" //中台目录
	ServerEtcdDir    = "/ServerList/"
	ConnectDir       = "/ConnectorServer/"
	FrontConnectDir  = ServerEtcdDir + FrontendDir + ConnectDir
	MiddleConnectDir = ServerEtcdDir + MiddlePlatDir + ConnectDir
)

var (
	APPEtcdDir              = "/APPCfgList/GameList/"
	MapFrontConnServer      map[string]ProxyNodeInfo
	MapMiddlePlatConnServer map[string]ProxyNodeInfo
	NewTCPAgent             func(*network.TCPConn) network.Agent
	ChanRoomFlag            chan bool
	ChanDataBase            chan DatabaseInfo
	ChanRedisInfo           chan RedisInfo
	ChanChildConf           chan EtcdChildConfig
	bInitialized            bool
)

type EtcdChildConfig struct {
	Action string
	Key    string
	Value  string
}

var roomInfoKey = ""

func StartEtcd() {
	MapFrontConnServer = make(map[string]ProxyNodeInfo)
	MapMiddlePlatConnServer = make(map[string]ProxyNodeInfo)
	ChanRoomFlag = make(chan bool, 1)
	ChanDataBase = make(chan DatabaseInfo, 100)
	ChanRedisInfo = make(chan RedisInfo, 100)
	ChanChildConf = make(chan EtcdChildConfig, 1000)
	writeRoomInfo2Etcd()
	//加载所需etcd配置
	if FrontentNode {
		Server.EtcdKey = append(Server.EtcdKey, FrontConnectDir)
	}
	if MiddlePlatNeed {
		Server.EtcdKey = append(Server.EtcdKey, MiddleConnectDir)
	}
	Server.EtcdKey = append(Server.EtcdKey, fmt.Sprintf("%s", DBEtcdDir))
	Server.EtcdKey = append(Server.EtcdKey, fmt.Sprintf("%s", RedisEtcdDir))
	if Server.IsGray == 1 {
		roomInfoKey = fmt.Sprintf("%s%s%s/Level_%d/%s",
			GrayPathPrefix, APPEtcdDir, Server.CfgDir, RoomInfo.RoomLevel, Server2Etcd.key)
		Server.EtcdKey = append(Server.EtcdKey, fmt.Sprintf("%s%s%s",
			GrayPathPrefix, APPEtcdDir, Server.CfgDir))
	} else {
		roomInfoKey = fmt.Sprintf("%s%s/Level_%d/%s",
			APPEtcdDir, Server.CfgDir, RoomInfo.RoomLevel, Server2Etcd.key)
		Server.EtcdKey = append(Server.EtcdKey, fmt.Sprintf("%s%s",
			APPEtcdDir, Server.CfgDir))
	}

	for _, v := range Server.EtcdKey {
		cfgetcd := client.Config{
			Endpoints: []string{Server.EtcdAddr},
			Transport: client.DefaultTransport,
			// set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: time.Second,
		}
		etcdClient, err := client.New(cfgetcd)
		if err != nil {
			log.Error("err:%v", err)
		}
		gameAPI := client.NewKeysAPI(etcdClient)
		strPath := Server.SysClusterName + "/" + v
		resp, err := gameAPI.Get(context.Background(), strPath,
			&client.GetOptions{Recursive: true, Sort: false, Quorum: true})
		if err != nil {
			log.Error("err:%v", err)
		}
		if resp != nil && resp.Node != nil {
			paraseEtcdNode(resp.Action, resp.Node)
		}
		go watchGateServer(strPath)
	}
	//注册本节点信息到etcd中心
	go registe2Etcd()
}

//创建静态节点信息,致后台管理
func writeRoomInfo2Etcd() {
	cfgetcd := client.Config{
		Endpoints: []string{Server.EtcdAddr},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	etcdClient, err := client.New(cfgetcd)
	if err != nil {
		log.Error("err:%v", err)
	}
	gameAPI := client.NewKeysAPI(etcdClient)
	roominfo, _ := json.Marshal(&RoomInfo)
	//灰度列表注册
	if Server.IsGray == 1 {
		_, err = gameAPI.Create(context.Background(), roomInfoKey, string(roominfo))
	} else {
		_, err = gameAPI.Create(context.Background(), roomInfoKey, string(roominfo))
	}

}
func registe2Etcd() {
	for {
		cfgetcd := client.Config{
			Endpoints: []string{Server.EtcdAddr},
			Transport: client.DefaultTransport,
			// set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: time.Second,
		}
		etcdClient, err := client.New(cfgetcd)
		if err != nil {
			log.Error("err:%v", err)
		}
		gameAPI := client.NewKeysAPI(etcdClient)
		Server2Etcd.value.CurOnlineNum = RoomInfo.CurOnlineNum
		strNode, _ := json.Marshal(&Server2Etcd.value)
		strValue := string(strNode[:])
		_, err = gameAPI.Set(context.Background(), fmt.Sprintf("%s/%s%s/Level_%d/%s",
			Server.SysClusterName, ServerEtcdDir, Server.CfgDir, RoomInfo.RoomLevel, Server2Etcd.key), strValue,
			&client.SetOptions{TTL: 10 * time.Second})
		if err != nil {
			log.Error("err:%v", err)
		}
		//log.Debug("resp:%v", resp)
		time.Sleep(9 * time.Second)
	}
}
func paraseEtcdNode(action string, node *client.Node) {
	if node.Dir == false {
		jsonConf2Struct(action, node.Key, node.Value)
	} else {
		for _, v := range node.Nodes {
			paraseEtcdNode(action, v)
		}
	}
}

func jsonConf2Struct(action string, key string, value string) {

	if strings.Contains(key, "ConnectorServer") {
		nodeinfo := ProxyNodeInfo{}
		json.Unmarshal([]byte(value), &nodeinfo)
		//如果是新增或者初次加载时，建立连接
		if action == "get" || action == "set" {
			if strings.Contains(key, FrontendDir) {
				if _, ok := MapFrontConnServer[nodeinfo.LocalAddr]; !ok {
					MapFrontConnServer[nodeinfo.LocalAddr] = nodeinfo
					ClientManager := &network.TCPClient{
						Addr:            nodeinfo.LocalAddr,
						ConnNum:         1,
						LittleEndian:    LittleEndian,
						PendingWriteNum: PendingWriteNum,
						AutoReconnect:   true,
						LenMsgLen:       LenMsgLen,
						MinMsgLen:       0,
						MaxMsgLen:       MaxMsgLen,
						NewAgent:        NewTCPAgent,
					}
					ClientManager.Start()
				}
			} else {
				if _, ok := MapMiddlePlatConnServer[nodeinfo.LocalAddr]; !ok {
					MapMiddlePlatConnServer[nodeinfo.LocalAddr] = nodeinfo
					ClientManager := &network.TCPClient{
						Addr:            nodeinfo.LocalAddr,
						ConnNum:         1,
						LittleEndian:    LittleEndian,
						PendingWriteNum: PendingWriteNum,
						AutoReconnect:   true,
						LenMsgLen:       LenMsgLen,
						MinMsgLen:       0,
						MaxMsgLen:       MaxMsgLen,
						NewAgent:        NewTCPAgent,
					}
					ClientManager.Start()
				}
			}

		}
	} else if key == roomInfoKey {
		roominfo := RoomInfoDef{}
		json.Unmarshal([]byte(value), &roominfo)
		RoomInfo.BasePoint = roominfo.BasePoint
		RoomInfo.CanWatchGame = roominfo.CanWatchGame
		// RoomInfo.GameStartCheckTime = roominfo.GameStartCheckTime
		// RoomInfo.GameStartPlayer = roominfo.GameStartPlayer
		// RoomInfo.MaxOnlineNum = roominfo.MaxOnlineNum
		// RoomInfo.MaxTableNum = roominfo.MaxTableNum
		// RoomInfo.MaxTableChair = roominfo.MaxTableChair
		RoomInfo.PropIdList = roominfo.PropIdList
		RoomInfo.RoomCoinDownLimit = roominfo.RoomCoinDownLimit
		RoomInfo.RoomCoinUpLimit = roominfo.RoomCoinUpLimit
		RoomInfo.ServerName = roominfo.ServerName
		RoomInfo.GameName = roominfo.GameName
		RoomInfo.ServerStatus = roominfo.ServerStatus
		RoomInfo.SitNoHandUpCheckTime = roominfo.SitNoHandUpCheckTime
		// RoomInfo.TaxRate = roominfo.TaxRate
		// RoomInfo.TableTax = roominfo.TableTax
		if bInitialized {
			if RoomInfo.ServerStatus == 0 {
				ChanRoomFlag <- false //房间关闭
			} else {
				ChanRoomFlag <- true //房间打开
			}
		}
		bInitialized = true
	} else if strings.Contains(key, "DBConf/") {
		dbinfo := DatabaseInfo{}
		json.Unmarshal([]byte(value), &dbinfo)
		bHas := false
		//有则更新
		for i, v := range Server.DbList {
			if v.Host == dbinfo.Host &&
				v.Port == dbinfo.Port &&
				v.DataBase == dbinfo.DataBase {
				Server.DbList[i] = dbinfo
				bHas = true
				break
			}
		}
		Server.MapDBName[dbinfo.DataBase] = true
		//无则创建,通知db
		if bInitialized {
			if !bHas {

				ChanDataBase <- dbinfo
			}
		} else {
			if !bHas {
				Server.DbList = append(Server.DbList, dbinfo)
			}
		}
	} else if strings.Contains(key, "RedisConf/") {
		dbinfo := RedisInfo{}
		json.Unmarshal([]byte(value), &dbinfo)
		bHas := false
		//有则更新
		for i, v := range Server.RedisList {
			if v.Addr == dbinfo.Addr &&
				v.RedisName == dbinfo.RedisName {
				Server.RedisList[i] = dbinfo
				bHas = true
				break
			}
		}
		Server.MapDBName[dbinfo.RedisName] = true
		//无则创建,通知db
		if bInitialized {
			if !bHas {
				ChanRedisInfo <- dbinfo
			}
		} else {
			if !bHas {
				Server.RedisList = append(Server.RedisList, dbinfo)
			}
		}
	} else {
		childConf := EtcdChildConfig{
			Action: action,
			Key:    key,
			Value:  value,
		}
		ChanChildConf <- childConf
	}
}
func watchGateServer(serverName string) {
	for {
		cfgetcd := client.Config{
			Endpoints: []string{Server.EtcdAddr},
			Transport: client.DefaultTransport,
			// set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: time.Second,
		}
		etcdClient, err := client.New(cfgetcd)
		if err != nil {
			log.Error("err:%v", err)
		}
		gameAPI := client.NewKeysAPI(etcdClient)
		rsp, err := gameAPI.Watcher(serverName, &client.WatcherOptions{Recursive: true}).Next(context.Background())
		if err == nil {
			paraseEtcdNode(rsp.Action, rsp.Node)
			//log.Debug("watchEtcdConf action:%s key:%s", rsp.Action, rsp.Node.Key)
		} else {
			time.Sleep(2 * time.Second)
		}
		gameAPI = nil
	}
}
