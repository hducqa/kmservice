package core

import (
	"github.com/go-xorm/xorm"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

const DefaultTag = "Center"
const DefaultInt = int64(0)

type PostTitle int

const (
	_ PostTitle = iota
	GET
	LINK
	UPDATE
	CONFIRM
	SUCCESS
	CONNECT
	FAILURE
	API_LIST
	EXCEPTION
	IS_ACTIVE
	FIND_LINK
	SUBSCRIBES
	LINK_SUBMIT
)

const (
	Log_Info  = "Info"
	Log_Warn  = "Warn"
	Log_Error = "Error"
)

const (
	NO_SUBSCRIBE_INFO            = "NO_SUBSCRIBE_INFO"
	KEY_NOT_EXIST                = "KEY_NOT_EXIST"
	DATA_LOCKED                  = "DATA_LOCKED"
	LINK_NOT_EXIST               = "LINK_NOT_EXIST"
	WITHOUT_PERMISSION           = "WITHOUT_PERMISSION"
	ORIGINAL_DATA_EXPIRED        = "ORIGINAL_DATA_EXPIRED"
	REQUEST_TYPE_EXCEPTION       = "REQUEST_TYPE_EXCEPTION"
	GET_DATA_FORM_EXECPTION      = "GET_DATA_FORM_EXECPTION"
	API_DATA_FORM_EXECPTION      = "API_DATA_FORM_EXECPTION"
	LINK_DATA_FORM_EXECPTION     = "LINK_DATA_FORM_EXECPTION"
	UPDATE_DATA_FORM_EXCEPTION   = "UPDATE_DATA_FORM_EXCEPTION"
	FINDLINK_DATA_FORM_EXECPTION = "FINDLINK_DATA_FORM_EXECPTION"
)

type ServiceState int

const (
	_ ServiceState = iota
	Stop
	Pending
	Active
)

type RegisterCenter struct {
	persistenceFilePath string                 //持久化文件路径
	DataMap             map[int64]interface{}  //共享文件库
	Subscribes          map[int64]Subscribe    //订阅名单
	sqlClient           *xorm.Engine           //数据库引擎
	ServiceCache        map[int64]MicroService //缓存所有服务基本信息
	ServiceActive       map[int64]ServiceState //记录服务是否活跃

	webSocketServer *socketio.Server //websocket服务
	logger          *logrus.Logger   //日志管理
	LogClient       *LogClient

	linkPool    map[string]LinkInfo
	socketPool  map[int64]net.Conn //TCP连接池
	connNum     int                //当前维护连接数
	maxPoolSize int                //最大连接数量

	persistenceChannel chan FileStorage   //数据存储通道
	updateChannel      chan UpdatePackage //数据更新通道
	rLocker            map[int64]bool     //读数据锁
}

type LogClient struct {
	SqlClient   *xorm.Engine
	ServiceId   int64
	ServiceName string
}

type Log struct {
	Id          int64     `json:"id"`
	ServiceId   int64     `json:"service_id"`
	ServiceName string    `json:"service_name"`
	Level       string    `json:"level"`
	File        string    `json:"file"`
	Line        int       `json:"line"`
	Message     string    `json:"message"`
	Time        time.Time `json:"time"`
}

type MicroService struct {
	Id           int64    `json:"id"`
	Name         string   `json:"name"`         //服务名称
	RootPath     string   `json:"rootPath"`     //服务所在根目录
	Language     string   `json:"language"`     //编码语言
	StartCommand string   `json:"startCommand"` //服务启动命令
	APIs         []API    `json:"APIs"`         //服务包含API内容
	IP           string   `json:"ip"`           //服务启动IP地址
	OwnerEmail   []string `json:"ownerEmail"`   //管理者邮箱
	Token        string   `json:"token"`        //服务密钥
}

type Subscribe struct {
	Id          int64   `json:"id"`
	Key         string  `json:"key"`
	Subscribers []int64 `json:"subscribers"`
	Writers     []int64 `json:"writers"`
	Description string  `json:"description"`
}

type API struct {
	Protocol    string `json:"protocol"` //API协议
	Route       string `json:"route"`    //路由
	RequestType string `json:"requestType"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ConnApply struct {
	Id    int64
	Token string
}

type DataGram struct {
	Tag       string
	ServiceId int64
	Data      Data
}

type UpdateRequset struct {
	Origin interface{}
	New    interface{}
}

type UpdatePackage struct {
	Tag       string
	ServiceId int64
	From      net.Conn
	Key       int64
	Request   UpdateRequset
}

type Data struct {
	Title     PostTitle
	Key       int64
	TimeStamp time.Time
	Body      interface{}
}

type FileStorage struct {
	DataMap map[int64]interface{} //共享文件库
}

type SubscribePrivilege struct {
	Read  bool
	Write bool
}

func (r RegisterCenter) PackageFile() FileStorage {
	return FileStorage{
		DataMap: r.DataMap,
	}
}

type LinkInfo struct {
	Key   string `json:"key"`
	Host  string `json:"host"`
	Port  string `json:"port"`
	Token string `json:"token"`
}

type LinkApply struct {
	Port string `json:"port"`
	Key  string `json:"key"`
}
