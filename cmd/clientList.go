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
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/murderxchip/aliyun-dts-dispatcher/define"
	"github.com/murderxchip/aliyun-dts-dispatcher/dts/subscribers"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"github.com/murderxchip/aliyun-dts-dispatcher/server"
	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/cobra"
)

// clientAddCmd represents the clientAdd command
var clientListCmd = &cobra.Command{
	Use:   "list",
	Short: "显示订阅列表",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := server.DtsRes.Etcd.Get(context.Background(), define.Prefix, clientv3.WithPrefix())
		if err != nil {
			log.Info(define.CmdTag, "load config error:", err)
			return
		}

		var cfg subscribers.SubscriberConfig
		for _, v := range resp.Kvs {
			err := json.Unmarshal(v.Value, &cfg)
			if err != nil {
				log.Info(define.CmdTag, err, string(v.Value))
				continue
			}
			cfgJson, _ := json.Marshal(cfg)
			cfgJsonStr := string(cfgJson)
			fmt.Printf("%s => %s\n", string(v.Key), cfgJsonStr)
		}

	},
}

func init() {
	clientCmd.AddCommand(clientListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
