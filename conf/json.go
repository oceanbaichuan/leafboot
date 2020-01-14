package conf

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/name5566/leaf/log"
)

var Server struct {
	LogLevel           string
	LogPath            string
	BuringPointLogPath string
	LogSplitHour       int32
	WSAddr             string
	CertFile           string
	KeyFile            string
	TCPAddr            string
	MaxConnNum         int32
	ConsolePort        int32
	ProfilePath        string
	PropIdList         []int32
	SysClusterName     string //本应用所属系统集群名称
	IsGray             int8   //是否灰度
	NodeName           string //gamename_gameid_roomlevel
	NodeID             string //pcid_roomserialid
	CfgDir             string
	CustomDBName       []string        //除平台公共库，自定义数据存储库名
	MapDBName          map[string]bool //除平台公共库，自定义数据存储库名
	//平台库
	DbList []DatabaseInfo
	//redis 列表
	RedisList []RedisInfo
	//etcd配置
	EtcdAddr string
	EtcdKey  []string
}

//DatabaseInfo 数据库部署采用M(应用主写)-M(应用主读)-S(后台读)模式
type DatabaseInfo struct {
	MinUID   int64 //分库用户ID起始
	MaxUID   int64 //分库用户ID截止
	Host     string
	Port     uint16
	UserName string
	Passwd   string
	DataBase string
	DbType   string // mysql,mssql等等
	IsMaster int8   //1:master 0：slave
	DbRWFlag string //write|read,readonly
}
type RedisInfo struct {
	Addr      string
	Passwd    string
	Slot      int32
	RedisName string
	MinUID    int64 //分库用户ID起始
	MaxUID    int64 //分库用户ID截止
}

//Node2EtcdInfo 写入etcd本身信息
type Node2EtcdInfo struct {
	NodeName     string //gamename_gameid_roomlevel
	RouterName   string //对外路由路径
	NodeID       string //pcid_roomserialid
	WSAddr       string
	TCPAddr      string
	CertFile     string
	KeyFile      string
	MaxConnNum   int32
	CurOnlineNum int32
}

var Server2Etcd struct {
	key   string
	value Node2EtcdInfo
}

//ProxyNodeInfo 代理服务器信息
type ProxyNodeInfo struct {
	NodeName   string //gamename_gameid_roomlevel
	NodeID     string //pcid_roomserialid
	LocalAddr  string
	RemoteAddr string
	CertFile   string
	KeyFile    string
	MaxConnNum int32
}

type RoomInfoDef struct {
	NodeName             string //gamename_gameid_roomlevel
	NodeID               string //pcid_roomserialid
	GameID               int32
	RoomLevel            int32
	MaxTableNum          int32 //最大桌子数
	MaxOnlineNum         int32 //最大在线人数
	MaxTableChair        int32 //一桌最多人数
	RoomCoinDownLimit    int64
	RoomCoinUpLimit      int64
	GameStartPlayer      int32 //满多少人开局
	GameStartCheckTime   int32 //检查开局间隔时长
	SitNoHandUpCheckTime int32 //入座未举手检测时长
	CanWatchGame         int8  //是否可以旁观游戏
	BasePoint            int32
	TaxRate              int32 //千分比
	TableTax             int32 //固定桌费
	ServerStatus         int32
	GameName             string
	ServerName           string
	PropIdList           []int32
	CurOnlineNum         int32
}

var RoomInfo RoomInfoDef

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Debug("%s", dir)
	//加载物理服配置
	data, err := ioutil.ReadFile("conf/server.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	Server.MapDBName = make(map[string]bool)
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
	if Server.SysClusterName == "" {
		log.Fatal("SysClusterName must not be empty string!!")
	}
	//解析数据库名称列表
	for _, v := range Server.DbList {
		Server.MapDBName[v.DataBase] = true
	}
	for _, v := range Server.RedisList {
		Server.MapDBName[v.RedisName] = true
	}
	for _, v := range Server.CustomDBName {
		Server.MapDBName[v] = true
	}
	roomflags := strings.Split(Server.NodeName, "_")
	if len(roomflags) <= 2 {
		Server.CfgDir = Server.NodeName
	} else {
		Server.CfgDir = strings.Join(roomflags[:2], "_")
	}
	if len(roomflags) > 1 {
		gameID, err := strconv.Atoi(roomflags[1])
		if err != nil {
			log.Fatal("%v", err)
		}
		if gameID > 0 {
			//加载房间配置
			data, err = ioutil.ReadFile("conf/room.json")
			if err != nil {
				log.Fatal("%v", err)
			}
			err = json.Unmarshal(data, &RoomInfo)
			log.Debug("RoomInfo:%v", RoomInfo)
			if err != nil {
				log.Fatal("%v", err)
			}
		}
		RoomInfo.GameID = int32(gameID)
	}
	if len(roomflags) > 2 {
		roomLevel, err := strconv.Atoi(roomflags[2])
		if err != nil {
			log.Fatal("%v", err)
		}
		RoomInfo.RoomLevel = int32(roomLevel)
	}
	RoomInfo.NodeName = Server.NodeName
	RoomInfo.NodeID = Server.NodeID
	nodeinfo := Node2EtcdInfo{}
	nodeinfo.CertFile = Server.CertFile
	nodeinfo.KeyFile = Server.KeyFile
	nodeinfo.MaxConnNum = Server.MaxConnNum
	nodeinfo.NodeID = Server.NodeID
	nodeinfo.NodeName = Server.NodeName
	nodeinfo.TCPAddr = Server.TCPAddr
	nodeinfo.WSAddr = Server.WSAddr
	nodeinfo.RouterName = Server.NodeName
	hashdig := md5.New()
	hashdig.Write([]byte(fmt.Sprintf("%s",
		Server.TCPAddr) + Server.NodeID))
	sercode := hashdig.Sum([]byte(""))
	Server2Etcd.key = fmt.Sprintf("%x", sercode)
	// strNode, err := json.Marshal(&nodeinfo)
	// if err != nil {
	// 	log.Fatal("%v", err)
	// }
	Server2Etcd.value = nodeinfo
}
