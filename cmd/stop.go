/*
Copyright © 2020 Romber Li <romber2001@gmail.com>

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
	"os"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"github.com/romberli/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/romberli/go-template/config"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop command",
	Long:  `stop the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
		)

		// init config
		err = initConfig()
		if err != nil {
			fmt.Println(fmt.Sprintf(config.ErrInitConfig, err.Error()))
		}

		// kill server with given pid
		if serverPid != constant.DefaultRandomInt {
			err = linux.KillServer(serverPid)
			if err != nil {
				log.Errorf(config.ErrKillServerWithPid, serverPid, err.Error())
				fmt.Println(fmt.Sprintf(config.ErrKillServerWithPid, serverPid))
			}

			log.Infof(config.InfoServerStop, serverPid, serverPidFile)
			fmt.Println(fmt.Sprintf(config.InfoServerStop, serverPid, serverPidFile))

			return
		}

		// get pid from pid file
		serverPidFile = viper.GetString(config.ServerPidFileKey)
		serverPid, err = linux.GetPidFromPidFile(serverPidFile)
		if err != nil {
			log.Errorf(config.ErrGetPidFromPidFile, serverPidFile, err.Error())
			fmt.Println(fmt.Sprintf(config.ErrGetPidFromPidFile, serverPidFile, err.Error()))
			os.Exit(constant.DefaultAbnormalExitCode)
		}

		// kill server with pid and pid file
		err = linux.KillServer(serverPid, serverPidFile)
		if err != nil {
			log.Errorf(config.ErrKillServerWithPidFile, serverPid, serverPidFile, err.Error())
			fmt.Println(fmt.Sprintf(config.ErrKillServerWithPidFile, serverPidFile, serverPid))
		}

		log.Infof(config.InfoServerStop, serverPid, serverPidFile)
		fmt.Println(fmt.Sprintf(config.InfoServerStop, serverPid, serverPidFile))
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")
	stopCmd.PersistentFlags().IntVar(&serverPid, "server-pid", constant.DefaultRandomInt, fmt.Sprintf("specify the server pid"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
