// Copyright © 2018 Job King'ori Maina <j@kingori.co>

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"

	"github.com/itskingori/sanaa/service"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print out the version of the software",
	Long:  "Print out the version of the software.",
	Run: func(cmd *cobra.Command, args []string) {
		version := service.GetVersion()

		fmt.Printf("Sanaa v%s\n", version.Str())
	},
}

// init initializes the command
func init() {
	RootCmd.AddCommand(versionCmd)
}
