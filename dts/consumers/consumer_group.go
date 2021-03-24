package consumers

import (
	"github.com/murderxchip/aliyun-dts-dispatcher/config"
	"github.com/murderxchip/aliyun-dts-dispatcher/define"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"sync"
)

type ConsumerGroup struct {
	consumers map[string]*Consumer
	m         sync.Mutex
	quit      int
	queue     chan DtsData
}

func NewConsumerGroup() (sg *ConsumerGroup) {
	sg = &ConsumerGroup{
		consumers: make(map[string]*Consumer),
		queue:     make(chan DtsData, 10000),
	}

	consumerGroupConfigs := config.Config().Dts
	for _, consumerGroupConfig := range consumerGroupConfigs {
		//load from config file
		cfg := ConsumerConfig{
			Name:     consumerGroupConfig.Name,
			Topic:    []string{consumerGroupConfig.Topic},
			User:     consumerGroupConfig.User,
			Password: consumerGroupConfig.Password,
			Broker:   []string{consumerGroupConfig.Broker},
			GroupId:  consumerGroupConfig.GroupId,
		}
		wg := sync.WaitGroup{}
		go func(consumerGroupConfig config.Dts, cfg ConsumerConfig) {
			wg.Add(1)
			consumer, err := NewConsumer(cfg, sg.queue)
			if err != nil {
				log.Warning(define.ConsumerTag, "consumer register error : "+consumerGroupConfig.Name)
			}
			sg.Register(consumer)
			wg.Done()
		}(consumerGroupConfig, cfg)
		wg.Wait()
	}
	return
}
func (s *ConsumerGroup) GetConsumer(name string) *Consumer {
	sub, exists := s.consumers[name]
	if !exists {
		return nil
	}

	return sub
}

func (s *ConsumerGroup) Register(consumer *Consumer) {
	s.m.Lock()
	defer s.m.Unlock()

	sub := s.GetConsumer(consumer.GetName())
	if sub != nil {
		//todo  stop
		(*sub).Shutdown()
	}
	s.consumers[consumer.GetName()] = consumer
	consumer.Start()
}

func (s *ConsumerGroup) Watch() {
	//todo update consumer config
}

func (s *ConsumerGroup) GetQueue() <-chan DtsData {
	return s.queue
}

func (s *ConsumerGroup) Shutdown() {
	//close consumers
	wg := sync.WaitGroup{}
	for _, c := range s.consumers {
		go func() {
			wg.Add(1)
			c.Shutdown()
			wg.Done()
		}()
	}
	wg.Wait()
	log.Info("ConsumerGroup", "shutdown")
	s.quit = 1
}

type Group map[string]Consumer

