// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// gotoCmd represents the goto command
var gotoCmd = &cobra.Command{
	Use:   "goto",
	Short: "Migrate to a specific version",
	Long:  `Migrate to a specific version`,
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl := cmd.Flags().Lookup("url").Value.String()
		if dbUrl == "" {
			env := cmd.Flags().Lookup("env").Value.String()
			if env != "" {
				dbUrl = viper.GetString(env)
			}

			if dbUrl == "" {
				fmt.Println("either `url` or `env` parameter must be provided.  `env` values can be stored in `.mg.yaml` or `$HOME/.mg.yaml`")
				return
			}
		}

		m, err := CreateMigration(dbUrl)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if len(args) < 1 {
			fmt.Println("down requires a specific number of steps to be applied")
			return
		}

		v, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = m.Migrate(uint(v))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(gotoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gotoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gotoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
