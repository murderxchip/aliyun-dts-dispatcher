package consumers

import "fmt"

type ConsumerConfig struct {
	Name     string
	User     string
	Password string
	Broker   []string
	GroupId  string
	Topic    []string
}

func NewConfig(broker, topic []string, groupid, user, password string) *ConsumerConfig {
	return &ConsumerConfig{
		User:     user,
		Password: password,
		Broker:   broker,
		GroupId:  groupid,
		Topic:    topic,
	}
}

func (c *ConsumerConfig) GetSASLUser() string {
	return fmt.Sprintf("%s-%s", c.User, c.GroupId)
}
