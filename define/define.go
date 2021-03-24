package define

/**
                   _ooOoo_
                  o8888888o
                  88" . "88
                  (| x_x |)
                  O\  0  /O
               ____/`---'\____
             .'  \\|     |//  `.
            /  \\|||  :  |||//  \
           /  _||||| -:- |||||-  \
           |   | \\\  -  /// |   |
           | \_|  ''\---/''  |   |
           \  .-\__  `-`  ___/-. /
         ___`. .'  /--.--\  `. . __
      ."" '<  `.___\_<|>_/___.'  >'"".
     | | :  `- \`.;`\ _ /`;.`/ - ` : | |
     \  \ `-.   \_ __\ /__ _/   .-` /  /
======`-.____`-.___\_____/___.-`____.-'======
                   `=---='
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
         佛祖保佑       永无BUG
*/

/**
1.config包中不可引用define文件
2.define不可引用其他包 ！！！！ 业务逻辑请放到别的包中实现 define包只放定义！！！！！！！！！！！！！！
3.log包不可引用其他包！！！！！！
*/

var GitSHA = ""
var GitDateTime = ""

const Version = "0.8.0-dev"

// tag
const (
	CmdTag         = "cmd"
	ConsumerTag    = "consumer"
	SubscribersTag = "subscriber"
	ServerTag      = "server"
	RocketMqHttp   = "rocker_mq_http"
	Config		   = "config"
	Kafka		   = "kafka"
)

// server env
const (
	ServerEnvProd  = "prod"
	ServerEnvTest  = "test"
	ServerEnvDev   = "dev"
	ServerEnvLocal = "local"
)

// etcd
const (
	Prefix = "subcfg"

	//TYPE
	SubscriberTypeRocketmqHttp = 1
	SubscriberTypeRocketmqTcp  = 2
	SubscriberTypeKafka011     = 3
)

var EtcdTypeMap = map[int]string{
	SubscriberTypeRocketmqHttp: "rocketMq",
	SubscriberTypeKafka011:     "kafka",
}

func TypeToName(etcdType int) string {
	if v, ok := EtcdTypeMap[etcdType]; ok {
		return v
	}
	return ""
}
