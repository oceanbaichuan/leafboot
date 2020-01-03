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
	LogSplitHour       int
	WSAddr             string
	CertFile           string
	KeyFile            string
	TCPAddr            string
	MaxConnNum         int
	ConsolePort        int
	ProfilePath        string
	PropIdList         []int
	//平台库
	DbList []DatabaseInfo
	//etcd配置
	EtcdAddr string
	EtcdKey  []string
}

type DatabaseInfo struct {
	MinUID   int64 //分库用户ID起始
	MaxUID   int64 //分库用户ID截止
	Host     string
	Port     uint16
	UserName string
	Passwd   string
	DataBase string
	DbType   string
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
	MaxConnNum   int
	CurOnlineNum int
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
	MaxConnNum int
}

type RoomInfoDef struct {
	NodeName             string //gamename_gameid_roomlevel
	NodeID               string //pcid_roomserialid
	CfgDir               string
	GameID               int
	RoomLevel            int
	MaxTableNum          int //最大桌子数
	MaxOnlineNum         int //最大在线人数
	MaxTableChair        int //一桌最多人数
	RoomCoinDownLimit    int64
	RoomCoinUpLimit      int64
	GameStartPlayer      int  //满多少人开局
	GameStartCheckTime   int  //检查开局间隔时长
	SitNoHandUpCheckTime int  //入座未举手检测时长
	CanWatchGame         int8 //是否可以旁观游戏
	BasePoint            int
	TaxRate              int //千分比
	TableTax             int //固定桌费
	ServerStatus         int
	GameName             string
	ServerName           string
	PropIdList           []int
	CurOnlineNum         int
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
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}

	//加载房间配置
	data, err = ioutil.ReadFile("conf/room.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &RoomInfo)
	if err != nil {
		log.Fatal("%v", err)
	}
	roomflags := strings.Split(RoomInfo.NodeName, "_")
	RoomInfo.CfgDir = strings.Join(roomflags[:2], "_")
	RoomInfo.GameID, err = strconv.Atoi(roomflags[1])
	if err != nil {
		log.Fatal("%v", err)
	}
	RoomInfo.RoomLevel, err = strconv.Atoi(roomflags[2])
	if err != nil {
		log.Fatal("%v", err)
	}
	nodeinfo := Node2EtcdInfo{}
	nodeinfo.CertFile = Server.CertFile
	nodeinfo.KeyFile = Server.KeyFile
	nodeinfo.MaxConnNum = Server.MaxConnNum
	nodeinfo.NodeID = RoomInfo.NodeID
	nodeinfo.NodeName = RoomInfo.NodeName
	nodeinfo.TCPAddr = Server.TCPAddr
	nodeinfo.WSAddr = Server.WSAddr
	nodeinfo.RouterName = RoomInfo.NodeName
	hashdig := md5.New()
	hashdig.Write([]byte(fmt.Sprintf("%s",
		Server.TCPAddr) + RoomInfo.NodeID))
	sercode := hashdig.Sum([]byte(""))
	Server2Etcd.key = fmt.Sprintf("%x", sercode)
	// strNode, err := json.Marshal(&nodeinfo)
	// if err != nil {
	// 	log.Fatal("%v", err)
	// }
	Server2Etcd.value = nodeinfo
}
