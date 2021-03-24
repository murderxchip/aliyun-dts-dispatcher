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
	"fmt"
	"github.com/murderxchip/aliyun-dts-dispatcher/config"
	"github.com/murderxchip/aliyun-dts-dispatcher/server"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func syncConfig()  {
	configChan := make(chan os.Signal, 1)
	signal.Notify(configChan, syscall.SIGUSR1)
	for {
		<-configChan
		config.ReloadConfig()
		fmt.Println("Reloaded config :" + config.Config().Version)
	}
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动服务",
	Long: `启动DTS服务`,
	Run: func(cmd *cobra.Command, args []string) {
		//go syncConfig()
		dts := server.NewDtsServer()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		s := <-c
		fmt.Println("Got signal:", s)
		dts.Shutdown()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
