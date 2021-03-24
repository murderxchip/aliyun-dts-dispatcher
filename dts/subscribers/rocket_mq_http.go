package subscribers

import (
	"encoding/json"
	"fmt"
	"github.com/murderxchip/aliyun-dts-dispatcher/define"
	"github.com/murderxchip/aliyun-dts-dispatcher/dts/consumers"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"github.com/aliyunmq/mq-http-go-sdk"
)

type RocketMqHttpSubscriber struct {
	SubscriberBase
	producer mq_http_sdk.MQProducer
}

func (s RocketMqHttpSubscriber) GetConfigHash() []byte {
	return s.ConfigHash
}

func (s RocketMqHttpSubscriber) GetForceUpdate() bool {
	return s.ForceUpdate
}

func NewRocketMqHttpSubscriber(config SubscriberRocketMqHttpConfig) *RocketMqHttpSubscriber {
	s := &RocketMqHttpSubscriber{
		SubscriberBase{
			Name:        config.Name,
			Consumer:    config.Consumer,
			Queue:       make(chan consumers.DtsData, 2000),
			ForceUpdate: config.ForceUpdate,
			ConfigHash:  config.Hash(),
		},
		nil,
	}

	client := mq_http_sdk.NewAliyunMQClient(config.Endpoint, config.AccessKey, config.SecretKey, "")
	log.Info(define.SubscribersTag, "client:", client)
	s.producer = client.GetProducer(config.InstanceId, config.Topic)
	s.SetTables(config.Tables)
	return s
}

func (s RocketMqHttpSubscriber) Run() {
	log.Info(define.SubscribersTag, "starting subscriber routine", s.Name)
	go func() {
		for {
			select {
			case data := <-s.Queue:
				//log.Info(define.SubscribersTag,"get dts message")
				str, err := json.Marshal(data)
				if err != nil {
					log.Error(define.SubscribersTag, "dts marshal error", err)
				} else {
					msg := mq_http_sdk.PublishMessageRequest{
						MessageBody: string(str),
						MessageTag:  data.Consumer,
						Properties:  nil,
						MessageKey:  "",
					}
					log.Debug(define.RocketMqHttp,"msg is",msg)

					ret, err := s.producer.PublishMessage(msg)
					log.Info(define.RocketMqHttp,"ret err",ret,err)
					if err != nil {
						//todo log error
						log.Error(define.SubscribersTag, "push msg error", err)
					} else {
						//successful
						log.Info(define.SubscribersTag, fmt.Sprintf("Publish ---->\n\tMessageId:%s, BodyMD5:%s, Body:%s\n", ret.MessageId, ret.MessageBodyMD5, msg.MessageBody))
					}
				}
			default:
				if s.Quit == 1 && len(s.Queue) == 0 {
					//todo quit subscriber
					log.Warning(define.SubscribersTag, "break sub routine")
					break
				}
			}
		}
	}()
}

