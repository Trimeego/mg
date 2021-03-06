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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// downCmd represents the down command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version",
	Long:  `Gets and prints the current version from the database specified`,
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
		v, dirty, err := m.Version()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		status := "(clean)"
		if dirty {
			status = "(dirty)"
		}
		fmt.Printf("Current Version: %d %s", v, status)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
