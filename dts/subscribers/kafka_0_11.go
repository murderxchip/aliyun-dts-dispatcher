package subscribers

import (
	"encoding/json"
	"fmt"
	"github.com/murderxchip/aliyun-dts-dispatcher/define"
	"github.com/murderxchip/aliyun-dts-dispatcher/dts/consumers"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"github.com/Shopify/sarama"
	"time"
)

type Kafka011Subscriber struct {
	SubscriberBase
	producer sarama.SyncProducer
	mqConfig MqConfig
}

func (s Kafka011Subscriber) GetConfigHash() []byte {
	return s.ConfigHash
}

func (s Kafka011Subscriber) GetForceUpdate() bool {
	return s.ForceUpdate
}

type MqConfig struct {
	Topics     []string `json:"topics"`
	Servers    []string `json:"servers"`
	ConsumerId string   `json:"consumerGroup"`
}

func NewKafka011Subscriber(config SubscriberKafka011Config) *Kafka011Subscriber {
	kafkaCfg := MqConfig{
		Topics:     config.Topics,
		Servers:    config.Servers,
		ConsumerId: config.ConsumerId,
	}
	s := &Kafka011Subscriber{
		SubscriberBase{
			Name:        config.Name,
			Consumer:    config.Consumer,
			Queue:       make(chan consumers.DtsData, 2000),
			ForceUpdate: config.ForceUpdate,
			ConfigHash:  config.Hash(),
		},
		nil,
		kafkaCfg,
	}
	s.SetTables(config.Tables)

	mqCfg := sarama.NewConfig()
	mqCfg.Producer.Return.Successes = true
	mqCfg.Version = sarama.V0_11_0_0

	if err := mqCfg.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka producer config invalidate. config: %v. err: %v", kafkaCfg, err)
		fmt.Println(msg)
		log.Fatal(define.Kafka,msg)
		panic(msg)
	}

	producer, err := sarama.NewSyncProducer(kafkaCfg.Servers, mqCfg)
	if err != nil {
		msg := fmt.Sprintf("Kafak producer create fail. err: %v", err)
		log.Fatal(define.Kafka,msg)
		panic(msg)
	}

	s.producer = producer
	return s
}

func (s Kafka011Subscriber) Run() {
	log.Info(define.SubscribersTag, "starting subscriber routine", s.Name)
	go func() {
		for {
			select {
			case data := <-s.Queue:
				//log.Info(define.SubscribersTag,"get dts message")
				str, err := json.Marshal(data)
				if err != nil {
					log.Error(define.SubscribersTag, "get dts error", err)
				} else {
					msg := &sarama.ProducerMessage{
						Topic:     s.mqConfig.Topics[0],
						Key:       sarama.StringEncoder(""),
						Value:     sarama.StringEncoder(str),
						Timestamp: time.Now(),
					}

					_, _, err := s.producer.SendMessage(msg)
					if err != nil {
						log.Error(define.SubscribersTag, "push msg error", err)
						return
					} else {
						//successful
						//log.Info(define.SubscribersTag, "Publish ok", str)
					}
				}
			default:
				if s.Quit == 1 && len(s.Queue) == 0 {
					//todo quit subscriber
					log.Error(define.SubscribersTag, "break sub routine")
					break
				}
				time.Sleep(time.Second * 1)
			}
		}
	}()
}
