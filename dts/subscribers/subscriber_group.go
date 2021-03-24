package subscribers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/murderxchip/aliyun-dts-dispatcher/define"
	"github.com/murderxchip/aliyun-dts-dispatcher/dts/consumers"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"sync"
)

type Map map[string]*Subscriber

type SubscriberGroup struct {
	subscribers map[string]Map
	m           sync.Mutex
	running     int
	quit        int
	queue       <-chan consumers.DtsData
	client      *clientv3.Client
	exit        chan interface{}
}

func NewSubscriberGroup(queue <-chan consumers.DtsData, client *clientv3.Client) (sg *SubscriberGroup) {
	sg = &SubscriberGroup{
		subscribers: make(map[string]Map),
		running:     0,
		queue:       queue,
		client:      client,
		exit:        make(chan interface{}, 1),
	}

	sg.LoadConfig()
	sg.Watch()

	log.Info(define.SubscribersTag, "subscriber group run")
	return
}

func (s *SubscriberGroup) LoadConfig() {
	//todo init subscribers
	log.Info(define.SubscribersTag, "load config start")
	resp, err := s.client.Get(context.Background(), define.Prefix, clientv3.WithPrefix())
	if err != nil {
		log.Error(define.SubscribersTag, "load config error:", err)
		return
	}

	log.Info(define.SubscribersTag, "loaded ", resp)
	for _, v := range resp.Kvs {
		s.SyncConfig(v.Value)
	}
	return
}

func (s *SubscriberGroup) Watch() {
	//todo watch config update from etcd
	ch := s.client.Watch(context.Background(), define.Prefix, clientv3.WithPrefix())

	go func() {
		for {
			select {
			case v := <-ch:
				var cfg SubscriberConfig
				for _, e := range v.Events {
					switch e.Type {
					case mvccpb.PUT:
						{
							log.Info(define.SubscribersTag, "watch config ", cfg)
							s.SyncConfig(e.Kv.Value)
						}
					case mvccpb.DELETE:
					}
				}
			}
		}
	}()
}

func (s *SubscriberGroup) SyncConfig(value []byte) {
	var cfg SubscriberConfig
	e := json.Unmarshal(value, &cfg)
	if e != nil {
		log.Error(define.SubscribersTag, "unmarshal error", e, string(value))
		return
	}

	log.Info(define.SubscribersTag, "loading subscriber config value:", string(value))
	//log.Info(define.SubscribersTag, "loading subscriber config:", cfg)

	switch cfg.Type {
	case define.SubscriberTypeRocketmqHttp:
		var cfgRocketMqHttp SubscriberRocketMqHttpConfig
		e := json.Unmarshal(value, &cfgRocketMqHttp)
		if e != nil {
			log.Error(define.SubscribersTag, "rocketmq http config unmarshal failed ", e)
		}
		log.Info(define.SubscribersTag, "cfgRocketMqHttp config:", cfgRocketMqHttp)
		sub := NewRocketMqHttpSubscriber(cfgRocketMqHttp)
		log.Debug(define.SubscribersTag, "sub config:", sub)
		s.Register(sub)
	case define.SubscriberTypeRocketmqTcp:
		//todo
	case define.SubscriberTypeKafka011:
		var cfgKafka011 SubscriberKafka011Config
		e := json.Unmarshal(value, &cfgKafka011)
		if e != nil {
			log.Error(define.SubscribersTag, "kafka config unmarshal failed ", e)
		}
		sub := NewKafka011Subscriber(cfgKafka011)
		s.Register(sub)
	}
}

func (s *SubscriberGroup) GetSubscriber(consumer, name string) *Subscriber {
	sub, exists := s.subscribers[consumer][name]
	if !exists {
		return nil
	}

	return sub
}

func (s *SubscriberGroup) Register(subscriber Subscriber) {
	s.m.Lock()
	defer s.m.Unlock()

	sub := s.GetSubscriber(subscriber.GetConsumerName(), subscriber.GetName())
	log.Debug(define.SubscribersTag,"sub is",sub)
	if sub != nil {
		if !subscriber.GetForceUpdate() {
			return
		}
		if !bytes.Equal(subscriber.GetConfigHash(), (*sub).GetConfigHash()) {
			log.Error(define.SubscribersTag,"confighash not equal")
			(*sub).Shutdown()
		}
	} else {
		_, ok := s.subscribers[subscriber.GetConsumerName()]
		if !ok {
			s.subscribers[subscriber.GetConsumerName()] = make(Map)
		}
	}

	s.subscribers[subscriber.GetConsumerName()][subscriber.GetName()] = &subscriber
	subscriber.Run()
}

func (s *SubscriberGroup) Dispatch() {
	s.m.Lock()
	defer s.m.Unlock()

	if s.running != 0 {
		return
	}

	s.running = 1
	log.Info(define.SubscribersTag, "subscriber group start")
	go func() {
		for {
			select {
			case data := <-s.queue:
				subs, ok := s.subscribers[data.Consumer]
				if ok {
					for _, v := range subs {
						if (*v).TableIn(data.Table) {
							log.Info(define.SubscribersTag, "consume ", data.Table,data)
							(*v).Consume(data)
						}
					}
				}else{
					log.Info(define.SubscribersTag, "subscriber group data consumer error not ok", ok)
				}
			default:
				if s.quit == 1 && len(s.queue) == 0 {
					//quit
					log.Info(define.SubscribersTag, "subscriber group quit.")
					s.exit <- true
					break
				}
				//time.Sleep(time.Second * 1)
			}
		}
	}()
}

func (s *SubscriberGroup) Shutdown() {
	//close subscribers
	wg := sync.WaitGroup{}
	for _, subMap := range s.subscribers {
		for _, v := range subMap {
			go func() {
				wg.Add(1)
				(*v).Shutdown()
				wg.Done()
			}()
		}
	}

	wg.Wait()
	//log.Info("SubscriberGroup", "shutdown")
	s.quit = 1
	//s.running = -1
}
