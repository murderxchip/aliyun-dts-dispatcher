package subscribers

import (
	"crypto/md5"
	"encoding/json"
)

type SubscriberConfig struct {
	Name     string `json:"name"`
	Consumer string `json:"consumer"`
	Type          int               `json:"type"`
	Tables        []string          `json:"tables"`
	ForceUpdate   bool              `json:"force_update"`
	TableHashRule map[string]string `json:"table_hash_rule"`
}

//type 1
type SubscriberRocketMqHttpConfig struct {
	SubscriberConfig
	Endpoint   string `json:"endpoint"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Topic      string `json:"topic"`
	InstanceId string `json:"instance_id"`
	GroupId    string `json:"group_id"`
}

//type 2
//todo

//type 3
type SubscriberKafka011Config struct {
	SubscriberConfig
	Servers    []string `json:"servers"`
	Topics     []string `json:"topics"`
	ConsumerId string   `json:"consumer_group"`
}

func (c *SubscriberConfig) Hash() []byte {
	data, _ := json.Marshal(*c)
	h := md5.New()
	return h.Sum(data)
}

func (c *SubscriberConfig) Check() bool {
	//todo check config validation

	return true
}
