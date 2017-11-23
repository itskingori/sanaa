// Copyright © 2017 Job King'ori Maina <j@kingori.co>
//
// This file is part of sanaa.
//
// sanaa is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// sanaa is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with sanaa. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"

	"github.com/itskingori/sanaa/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the application web server",
	Long:  `Start the application web server listening on the configured binding address and port.`,
	Args: func(cmd *cobra.Command, args []string) error {
		rt, _ := cmd.Flags().GetInt("request-ttl")

		if rt < service.MinRequestTTL {
			return fmt.Errorf("set request-ttl is %d, yet the minimum is %d", rt, service.MinRequestTTL)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting the server")

		client := service.NewClient()
		client.StartServer()
	},
}

// init initializes the command
func init() {
	RootCmd.AddCommand(serverCmd)

	// Add flags to serverCmd
	serverCmd.PersistentFlags().String("binding-address", "0.0.0.0", "address to bind to and listen for requests")
	serverCmd.PersistentFlags().Int("binding-port", 8080, "port to bind to and listen for requests")
	serverCmd.PersistentFlags().Int("request-ttl", 86400, "persistence window of data in ephemeral storage, in seconds")

	// Bind serverCmd flags with viper configuration
	viper.BindPFlag("server.binding_address", serverCmd.PersistentFlags().Lookup("binding-address"))
	viper.BindPFlag("server.binding_port", serverCmd.PersistentFlags().Lookup("binding-port"))
	viper.BindPFlag("server.request_ttl", serverCmd.PersistentFlags().Lookup("request-ttl"))
}
