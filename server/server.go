package server

import (
	"fmt"
	"github.com/murderxchip/aliyun-dts-dispatcher/config"
	"github.com/murderxchip/aliyun-dts-dispatcher/define"
	"github.com/murderxchip/aliyun-dts-dispatcher/dts/consumers"
	"github.com/murderxchip/aliyun-dts-dispatcher/dts/subscribers"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
)

type DtsResource struct {
	Etcd *clientv3.Client
}

func (d *DtsResource) Close() {
	d.Etcd.Close()
}

var DtsRes DtsResource

type DtsServer struct {
	name            string
	http            *gin.Engine
	consumerGroup   *consumers.ConsumerGroup
	subscriberGroup *subscribers.SubscriberGroup
	//viper           *viper.Viper
	//etcd            *clientv3.Client
	exit chan int
}

var Server *DtsServer

func NewDtsServer() *DtsServer {
	//Server = DtsServer{}
	if Server != nil {
		return Server
	}

	Server = &DtsServer{}

	log.Info(define.ServerTag, "starting consumers")
	cg := consumers.NewConsumerGroup()
	log.Info(define.ServerTag, "starting subscribers")

	sg := subscribers.NewSubscriberGroup(cg.GetQueue(), DtsRes.Etcd)
	sg.Dispatch()

	Server.subscriberGroup = sg
	Server.consumerGroup = cg
	Server.name = "Global Dts Server"

	//init gin
	Server.http = gin.Default()
	Server.http.GET("/version", Version)
	Server.http.GET("/client/list", ClientList)
	Server.http.POST("/client/del", ClientDel)
	Server.http.POST("/client/put", ClientPut)
	Server.http.POST("/client/test",Test)

	go func() {
		log.Info(define.ServerTag, "http starting")
		err := Server.http.Run(fmt.Sprintf(":%d", config.Config().Http.Port))
		if err != nil {
			log.Error(define.ServerTag, "http start failed:", err)
		}
	}()

	log.Info(define.ServerTag, "server starting")
	return Server
	//init
}

func (s *DtsServer) GetName() string {
	return s.name
}

func (s *DtsServer) Shutdown() {
	//shutdown modules
	log.Info(define.ServerTag, "server shutting down")

	s.consumerGroup.Shutdown()
	s.subscriberGroup.Shutdown()
	DtsRes.Close()

	log.Info(define.ServerTag, "done")
}
