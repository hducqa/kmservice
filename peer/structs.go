package peer

import (
	"github.com/go-xorm/xorm"
	"github.com/hducqa/kmservice/core"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type LinkType string

const (
	STOP     = "STOP"
	START    = "START"
	CUSTOM   = "CUSTOM"
	CONFIRM  = "CONFIRM"
	SUCCESS  = "SUCCESS"
	TRANSFER = "TRANSFER"
)

type Peer struct {
	centerIP    string
	centerPort  string
	token       string
	ServiceId   int64
	ServiceName string

	peerData          map[int64]interface{}
	getList           map[int64]bool
	updateRequestList map[int64]int //订阅更新申请状态，1为申请中，2为申请成功，-1为申请失败,0为可申请
	subscribeKeys     map[string]int64

	pendingList    map[string]PendingGram //等待队列
	pendingChannel chan PendingChannelItem

	LinkInfos map[string]core.LinkInfo //连接配置

	Links map[string]*Link

	logger    *logrus.Logger
	LogClient *core.LogClient
	sqlClient *xorm.Engine

	maxErrorTimes int
	connection    net.Conn
	errorTimes    int
	filePath      string

	readChannel chan byte
	gramChannel chan core.DataGram
}

type PeerConfig struct {
	CenterIP    string `json:"center_ip"`
	CenterPort  string `json:"center_port"`
	Token       string `json:"token"`
	ServiceId   int64  `json:"service_id"`
	ServiceName string `json:"service_name"`
	FilePath    string `json:"peer_file_path"`
}

type PendingChannelItem struct {
	Delete bool
	Tag    string
	Item   PendingGram
}

type PendingGram struct {
	Time        time.Time
	ResendTimes int
	Message     core.DataGram
}

type LinkInfo struct {
	port  string
	token string
}

type LinkGram struct {
	Tag       string
	Type      LinkType
	CustomKey string
	Body      interface{}
}

type LinkField struct {
	stop          bool
	conn          net.Conn
	DataChannel   chan interface{}
	CustomChannel chan LinkGram
	logger        *logrus.Logger
	logClient     *core.LogClient

	readChannel chan byte
	gramChannel chan LinkGram

	pendingList    map[string]PendingLinkGram
	pendingChannel chan PendingLinkChannelItem
}

type PendingLinkGram struct {
	linkGram    LinkGram
	resendTimes int
	Time        time.Time
}

type PendingLinkChannelItem struct {
	Delete bool
	Tag    string
	Item   PendingLinkGram
}

type Link struct {
	logger     *logrus.Logger
	logClient  *core.LogClient
	Token      string
	LinkNumber int
	LinkFields []LinkField
	DataField  []interface{}
}

type LinkApply struct {
	Token string `json:"token"`
	Desc  string `json:"desc"`
}
