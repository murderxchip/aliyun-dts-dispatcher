/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"github.com/murderxchip/aliyun-dts-dispatcher/cmd"
	"github.com/murderxchip/aliyun-dts-dispatcher/config"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"github.com/murderxchip/aliyun-dts-dispatcher/server"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func init() {
	log.NewLogConfig()
	config.ConfigServer = config.Config()

	fmt.Println(config.ConfigServer.Etcd)
	etcdCfg := clientv3.Config{
		Endpoints:   config.ConfigServer.Etcd.Endpoints,
		Username:    config.ConfigServer.Etcd.Username,
		Password:    config.ConfigServer.Etcd.Password,
		DialTimeout: time.Second * 10,
	}
	etcdObj, err := clientv3.New(etcdCfg)
	if err != nil {
		fmt.Println(etcdCfg)
		panic("ectd init failed")
	}

	server.DtsRes.Etcd = etcdObj
}

func main() {
	defer log.Logger.Close()
	cmd.Execute()
}
