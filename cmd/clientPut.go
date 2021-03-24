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
	"encoding/json"
	"fmt"
	"github.com/murderxchip/aliyun-dts-dispatcher/define"
	"github.com/murderxchip/aliyun-dts-dispatcher/dts/subscribers"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"github.com/spf13/cobra"
)

// clientAddCmd represents the clientAdd command
var clientAddCmd = &cobra.Command{
	Use:   "put",
	Short: "推送配置",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		//if len(args) == 3 && args[0] == "-f" {
		//
		//} else if len(args) == 2 {
		//
		//} else {
		//	fmt.Println("输入错误 ")
		//}
		var cfg subscribers.SubscriberConfig
		err := json.Unmarshal([]byte(args[1]), &cfg)
		if err != nil {
			log.Error(define.CmdTag, "")
			return
		}
		//if cfg.Name != "" && cfg.Consumer != "" && len(cfg.Endpoints) > 0 && cfg.AccessKey != "" && cfg.SecretKey != "" &&
		//	len(cfg.Topics) > 0 && cfg.InstanceId != "" && len(cfg.Tables) > 0 && cfg.Type > 0 && cfg.Type < 4 {
		//	log.Error(define.CmdTag,"")
		//	return
		//}


		//resp, err := server.DtsRes.Etcd.Put()
		fmt.Println(args[0], args[1])
	},
}

func init() {
	clientCmd.AddCommand(clientAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
