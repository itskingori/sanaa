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

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the application background worker",
	Long:  `Start the application background worker to process enqeueud jobs.`,
	Args: func(cmd *cobra.Command, args []string) error {
		cv, _ := cmd.Flags().GetInt("concurrency")

		if cv < service.MinWorkerConcurrency {
			return fmt.Errorf("set concurrency is %d, yet the minimum is %d", cv, service.MinWorkerConcurrency)
		}

		if cv > service.MaxWorkerConcurrency {
			return fmt.Errorf("set concurrency is %d, yet the maximum is %d", cv, service.MaxWorkerConcurrency)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting the worker")

		client := service.NewClient()
		client.StartWorker()
	},
}

// init initializes the command
func init() {
	RootCmd.AddCommand(workerCmd)

	// Add flags to workerCmd
	workerCmd.PersistentFlags().Int("concurrency", 2, "number of conversion jobs that can be processed at a time, maximum is 10")
	workerCmd.PersistentFlags().String("s3-bucket", "", "the name of the S3 bucket to use when storing rendered files ")
	workerCmd.PersistentFlags().String("s3-region", "", "the region of the S3 bucket to use when storing rendered files")

	// Bind workerCmd flags with viper configuration
	viper.BindPFlag("worker.concurrency", workerCmd.PersistentFlags().Lookup("concurrency"))
	viper.BindPFlag("worker.s3_bucket", workerCmd.PersistentFlags().Lookup("s3-bucket"))
	viper.BindPFlag("worker.s3_region", workerCmd.PersistentFlags().Lookup("s3-region"))
}
