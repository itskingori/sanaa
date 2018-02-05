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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var cfgFile string
var verbose bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sanaa",
	Short: "A HTML to PDF/Image conversion HTTP API",
	Long:  `A HTML to PDF/Image conversion HTTP API powered by wkhtmltopdf and wkhtmltoimage.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("Failed executing command, %s\n", err)
		os.Exit(-1)
	}
}

// init initializes the command
func init() {
	cobra.OnInitialize(initConfig)

	// Add flags to RootCmd
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output for debugging")
	RootCmd.PersistentFlags().String("redis-host", "127.0.0.1", "host of redis server")
	RootCmd.PersistentFlags().Int("redis-port", 6379, "port of redis server")
	RootCmd.PersistentFlags().String("redis-namespace", "sanaa", "namespace to use when storing data in redis server")

	// Bind RootCmd flags with viper configuration
	viper.BindPFlag("redis.host", RootCmd.PersistentFlags().Lookup("redis-host"))
	viper.BindPFlag("redis.port", RootCmd.PersistentFlags().Lookup("redis-port"))
	viper.BindPFlag("redis.namespace", RootCmd.PersistentFlags().Lookup("redis-namespace"))
}

// initConfig applies initial configuration
func initConfig() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
