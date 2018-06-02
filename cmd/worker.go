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
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the application background worker",
	Long:  `Start the application background worker to process enqeueud jobs.`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.NoArgs(cmd, args)

		err = validateWorkerConcurrency(cmd)
		if err != nil {

			return err
		}

		err = validateWorkerMaxRetries(cmd)
		if err != nil {

			return err
		}

		err = validateWorkerS3Bucket(cmd)
		if err != nil {

			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("starting the worker")

		client := service.NewClient()
		client.StartWorker()
	},
}

// init initializes the command
func init() {
	RootCmd.AddCommand(workerCmd)

	// Add flags to workerCmd
	workerCmd.PersistentFlags().Int("concurrency", 2, "number of conversion jobs that can be processed at a time, maximum is 10")
	workerCmd.PersistentFlags().Int("max-retries", 1, "maximum number of times to retry a job on failure")
	workerCmd.PersistentFlags().String("s3-bucket", "", "the name of the S3 bucket to use when storing rendered files ")

	// Bind workerCmd flags with viper configuration
	viper.BindPFlag("worker.concurrency", workerCmd.PersistentFlags().Lookup("concurrency"))
	viper.BindPFlag("worker.max-retries", workerCmd.PersistentFlags().Lookup("max-retries"))
	viper.BindPFlag("worker.s3_bucket", workerCmd.PersistentFlags().Lookup("s3-bucket"))
}

// validateWorkerConcurrency validate the concurrency flag
func validateWorkerConcurrency(cmd *cobra.Command) error {
	cv, _ := cmd.Flags().GetInt("concurrency")

	if cv < service.MinWorkerConcurrency {
		return fmt.Errorf("set concurrency is %d, yet the minimum is %d", cv, service.MinWorkerConcurrency)
	}

	if cv > service.MaxWorkerConcurrency {
		return fmt.Errorf("set concurrency is %d, yet the maximum is %d", cv, service.MaxWorkerConcurrency)
	}

	return nil
}

// validateWorkerMaxRetries validate the max-retries flag
func validateWorkerMaxRetries(cmd *cobra.Command) error {
	mrv, _ := cmd.Flags().GetInt("max-retries")

	if mrv < service.MinWorkerMaxRetries {
		return fmt.Errorf("set max-retries is %d, yet the minimum is %d", mrv, service.MinWorkerMaxRetries)
	}

	return nil
}

// validateWorkerS3Bucket validate the s3-bucket flag
func validateWorkerS3Bucket(cmd *cobra.Command) error {
	sbv, _ := cmd.Flags().GetString("s3-bucket")

	if sbv == "" {
		return fmt.Errorf("the S3 bucket name cannot be empty, set --s3-bucket")
	}

	return nil
}
