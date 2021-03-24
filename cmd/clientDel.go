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
	"errors"
	"fmt"
	"github.com/murderxchip/aliyun-dts-dispatcher/server"
	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/cobra"
	"time"
)

// clientDelCmd represents the clientDel command
var clientDelCmd = &cobra.Command{
	Use:   "del",
	Short: "删除订阅",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var response *clientv3.GetResponse
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		ch := make(chan int, 1)

		go func() {
			response, err = server.DtsRes.Etcd.Get(ctx, args[0], clientv3.WithKeysOnly(), clientv3.WithLimit(1))
			if err != nil {
				ch <- 1
			} else {
				if len(response.Kvs) == 0 {
					err = errors.New("could not found key " + args[0])
				}
				server.DtsRes.Etcd.Delete(ctx, args[0])
				ch <- 1
			}
		}()

		select {
		case <-ch:
			if err != nil {
				fmt.Println("del error: ", err)
			} else {

				fmt.Println("del success")
			}
		case <-ctx.Done():
			fmt.Printf("del error: %s", err)
		}

	},
}

func init() {
	clientCmd.AddCommand(clientDelCmd)

	clientDelCmd.PersistentFlags().String("name", "", "根据订阅名称删除订阅配置")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientDelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientDelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
