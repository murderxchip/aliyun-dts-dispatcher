package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/murderxchip/aliyun-dts-dispatcher/define"
	"github.com/murderxchip/aliyun-dts-dispatcher/dts/subscribers"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ClientDelHeader struct {
	Name string `json:"name"`
}

type ClientPutHeader struct {
	Config string `json:"config"`
}

func Version(c *gin.Context) {
	res := gin.H{
		"gin":          gin.Version,
		"server":       define.Version,
		"git-sha":      define.GitSHA,
		"git-datetime": define.GitDateTime,
	}
	c.JSON(http.StatusOK, res)
	log.Info(define.ServerTag,"version res",res)
}

func ClientDel(c *gin.Context) {
	var err error
	header := ClientDelHeader{}
	err = c.ShouldBindJSON(&header)
	if err != nil{
		log.Error(define.ServerTag,"clientDel params",header)
	}
	name := header.Name
	//name := c.PostForm("name")
	log.Info(define.ServerTag,"client params",name)

	var response *clientv3.GetResponse
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ch := make(chan int, 1)

	go func() {
		response, err = DtsRes.Etcd.Get(ctx, name, clientv3.WithKeysOnly(), clientv3.WithLimit(1))
		if err != nil {
			ch <- 1
		} else {
			if len(response.Kvs) == 0 {
				err = errors.New("could not found key " + name)
			}
			DtsRes.Etcd.Delete(ctx, name)
			ch <- 1
		}
	}()

	select {
	case <-ch:
		if err != nil {
			c.JSON(http.StatusOK, BuildResponse(define.StatFail, fmt.Sprintf("删除失败: %s", err)))
		} else {
			c.JSON(http.StatusOK, BuildResponse(define.StatSuccess, "删除成功"))
		}
	case <-ctx.Done():
		c.JSON(http.StatusOK, BuildResponse(define.StatFail, fmt.Sprintf("删除超时: %s", err)))
	}
}

func ClientPut(c *gin.Context) {
	header := ClientPutHeader{}
	err := c.ShouldBindJSON(&header)
	if err != nil{
		log.Error(define.ServerTag,"ClientPut params",header)
	}
	config := header.Config
	//config := c.PostForm("config")
	log.Info(define.ServerTag,"ClientPut params",config)

	if config == "" {
		c.JSON(http.StatusOK, BuildResponse(define.StatFail, "config is empty"))
	}

	var cfg subscribers.SubscriberConfig
	err = json.Unmarshal([]byte(config), &cfg)
	if err != nil {
		c.JSON(http.StatusOK, BuildResponse(define.StatFail, fmt.Sprintf("config format error:%s", err)))
		return
	}

	if cfg.Type < 1 || cfg.Type > 3 || cfg.Name == "" {
		c.JSON(http.StatusOK, BuildResponse(define.StatFail, "type between 1 AND 3 name not empty"))
		return
	}

	key := define.Prefix + "-" + cfg.Name + "-" + define.TypeToName(cfg.Type)
	_, err = DtsRes.Etcd.Put(context.Background(), key, config)
	if err != nil {
		c.JSON(http.StatusOK, BuildResponse(define.StatFail, err))
		return
	}
	switch cfg.Type {
	case define.SubscriberTypeRocketmqHttp:
		var cfgRocketMqHttp subscribers.SubscriberRocketMqHttpConfig
		err := json.Unmarshal([]byte(config), &cfgRocketMqHttp)
		if err != nil {
			c.JSON(http.StatusOK, BuildResponse(define.StatFail, err))
			return
		}
		c.JSON(http.StatusOK, BuildResponse(define.StatSuccess, cfgRocketMqHttp))
		log.Info(define.ServerTag,"ClientPut res",cfgRocketMqHttp)
	case define.SubscriberTypeKafka011:
		var cfgKafka011 subscribers.SubscriberKafka011Config
		err := json.Unmarshal([]byte(config), &cfgKafka011)
		if err != nil {
			c.JSON(http.StatusOK, BuildResponse(define.StatFail, err))
			return
		}
		c.JSON(http.StatusOK, BuildResponse(define.StatSuccess, cfgKafka011))
		log.Info(define.ServerTag,"ClientPut res",cfgKafka011)
	}
	return
}

func ClientList(c *gin.Context) {
	resp, err := DtsRes.Etcd.Get(context.Background(), define.Prefix, clientv3.WithPrefix())
	if err != nil {
		c.JSON(http.StatusOK, BuildResponse(define.StatFail, fmt.Sprintf("load config error:%s\n", err)))
		return
	}

	item := make(map[string]interface{})
	for _, v := range resp.Kvs {
		var cfg subscribers.SubscriberConfig
		err := json.Unmarshal([]byte(v.Value), &cfg)
		if err != nil {
			c.JSON(http.StatusOK, BuildResponse(define.StatFail, err))
			return
		}
		switch cfg.Type {
		case define.SubscriberTypeRocketmqHttp:
			var cfgRocketMqHttp subscribers.SubscriberRocketMqHttpConfig
			err := json.Unmarshal([]byte(v.Value), &cfgRocketMqHttp)
			if err != nil {
				c.JSON(http.StatusOK, BuildResponse(define.StatFail, err))
				return
			}
			item[string(v.Key)] = cfgRocketMqHttp
		case define.SubscriberTypeKafka011:
			var cfgKafka011 subscribers.SubscriberKafka011Config
			err := json.Unmarshal([]byte(v.Value), &cfgKafka011)
			if err != nil {
				c.JSON(http.StatusOK, BuildResponse(define.StatFail, err))
				return
			}
			item[string(v.Key)] = cfgKafka011
		}

	}
	c.JSON(http.StatusOK, BuildResponse(define.StatSuccess, item))
	log.Info(define.ServerTag,"ClientList res",item)
}

func Test(c *gin.Context) {
	header := ClientPutHeader{}
	err := c.ShouldBindJSON(&header)

	if err != nil{
		log.Error(define.ServerTag,"ClientPut params",header,err)
	}
	fmt.Println(header)
	config := header.Config
	fmt.Println(config)
	c.JSON(http.StatusOK, BuildResponse(define.StatSuccess, config))
}
