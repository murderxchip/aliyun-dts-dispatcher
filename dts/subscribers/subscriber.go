package subscribers

import (
	"fmt"
	"github.com/murderxchip/aliyun-dts-dispatcher/define"
	"github.com/murderxchip/aliyun-dts-dispatcher/dts/consumers"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"github.com/murderxchip/aliyun-dts-dispatcher/tools"
	"time"
)

type Subscriber interface {
	//订阅名称
	GetName() string
	//订阅数据库名称
	GetConsumerName() string
	//监听数据表列表
	GetTables() map[string]int
	//数据表是否订阅
	TableIn(table string) bool
	//消费dts数据包
	Consume(data consumers.DtsData)
	//启动线程
	Run()
	//退出清理
	Shutdown()
	//更新配置是否重启服务
	GetForceUpdate() bool
	//
	GetConfigHash() []byte
}

type SubscriberBase struct {
	//订阅名称
	Name string
	//数据源名称
	Consumer string
	//消费队列
	Queue chan consumers.DtsData
	//结束
	Quit int
	//是否已停止
	exit chan interface{}
	//
	Tables map[string]int
	//
	ForceUpdate bool
	//
	ConfigHash []byte
}

func (sb *SubscriberBase) WaitExit() {
	tools.WaitTimeout(fmt.Sprintf("subscriber: %s", sb.GetName()), sb.exit, 5*time.Second)
}

func (sb *SubscriberBase) Consume(data consumers.DtsData) {
	sb.Queue <- data
}

func (sb *SubscriberBase) Shutdown() {
	sb.Quit = 1
	sb.WaitExit()
}

func (sb *SubscriberBase) GetName() string {
	return sb.Name
}

func (sb *SubscriberBase) GetConsumerName() string {
	return sb.Consumer
}

func (sb *SubscriberBase) SetTables(tables []string) {
	tbl := make(map[string]int)
	for _, v := range tables {
		tbl[v] = 1
	}
	sb.Tables = tbl
	log.Debug(define.SubscribersTag,"set tables:", sb.Tables)
}

func (sb *SubscriberBase) GetTables() map[string]int {
	return sb.Tables
}

func (sb *SubscriberBase) TableIn(table string) bool {
	_, exist := sb.Tables[table]
	return exist
}
