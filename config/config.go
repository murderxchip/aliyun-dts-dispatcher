package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"path/filepath"
	"sync"
)

/**
只允许外部调用config，不允许config中调用别的包！！！！！！！！！！
防止循环引用！！！！！！！！！！
*/

var ConfigServer *tomlConfig

type tomlConfig struct {
	Version string
	Env     string
	Dts     []Dts
	Etcd    Etcd
	Log     Log
	Http    Http
}

type Etcd struct {
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}

type Dts struct {
	Name     string `json:"name"`
	Broker   string `json:"broker"`
	User     string `json:"user"`
	Password string `json:"password"`
	GroupId  string `json:"groupid"`
	Topic    string `json:"topic"`
}

//type Subscriber struct {
//	Name      string   `json:"name"`
//	Servers   []string `json:"server"`
//	GroupId   string   `json:"groupid"`
//	Topic    []string `json:"topic"`
//	AccessKey string   `json:"accesskey"`
//	Secretkey string   `json:"secretkey"`
//}

//日志配置文件
type Log struct {
	Level       string `json:"level"`
	LogPath     string `json:"logPath"`
	RotateDaily bool   `json:"rotateDaily"`
	Rotate      bool   `json:"rotate"`
}

type Http struct {
	Port int `json:"port"`
}

//func init() {
//	addFlag(flag.CommandLine)
//}
//
//func addFlag(fs *flag.FlagSet) {
//	fs.StringVar(&configFile, "config", "", "config file.")
//}

var (
	Cfg     *tomlConfig
	once    sync.Once
	cfgLock = new(sync.RWMutex)
)

func Config() *tomlConfig {
	once.Do(ReloadConfig)
	//cfgLock.RLock()
	//defer cfgLock.RUnlock()
	return Cfg
}

//reload配置文件
func ReloadConfig() {
	cfgLock.Lock()
	defer cfgLock.Unlock()
	filePath, err := filepath.Abs("./config/config.toml")
	if err != nil {
		fmt.Println("this is config, err:89")
		panic(err)
		//log.Error(define.Config,err)
	}
	config := new(tomlConfig)
	_, err = toml.DecodeFile(filePath, &config)
	if err != nil {
		fmt.Println("this is config, err:96")
		panic(err)
		//log.Error(define.Config,err)
	}
	Cfg = config
}
