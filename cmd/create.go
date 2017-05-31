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
	"io/ioutil"
	"math"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new migration scripts",
	Long:  `Creates a new set og migrations scripts that follow the proper naming convention.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		ts := math.Floor(float64(time.Now().UnixNano()) / (1000000000))
		if len(args) == 0 {
			fmt.Println("An argument must be provided.  This should be a snake case description eg. 'create_portal_user', of the script and will be used in the file creation.")
			return
		}

		upFile := fmt.Sprintf("%.0f_%s.up.sql", ts, args[0])
		downFile := fmt.Sprintf("%.0f_%s.down.sql", ts, args[0])

		err := ioutil.WriteFile(upFile, []byte("-- Migration Up\n\n"), 0644)
		if err != nil {
			log.Errorln(err.Error())
			return
		}

		err = ioutil.WriteFile(downFile, []byte("-- Migration Down\n\n"), 0644)
		if err != nil {
			log.Errorln(err.Error())
			return
		}

	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

}
