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

package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/gocraft/work"
	"github.com/itskingori/go-wkhtml/wkhtmltox"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

const (
	conversionQueue = "convert"

	// MinWorkerConcurrency is the minimum number of concurrent jobs the worker
	// should process
	MinWorkerConcurrency = 1

	// MaxWorkerConcurrency is the maximum number of concurrent jobs the worker
	// should process
	MaxWorkerConcurrency = 10

	// MinWorkerMaxRetries is minimum number that can be set for the worker's
	// maximum-retries
	MinWorkerMaxRetries = 0
)

type workerContext struct {
	client Client
}

func (ctx *workerContext) convert(job *work.Job) error {
	cl := NewClient()
	conn := cl.redisPool.Get()
	defer conn.Close()

	// Extract job parameter i.e. UUID
	jid := job.ArgString("uuid")
	log.WithFields(log.Fields{
		"uuid": jid,
	}).Info("picked up conversion job from queue")

	// Fetch all the job details
	cj, _, err := cl.fetchConversionJob(jid)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": jid,
		}).Errorf("error: %v", err)

		return err
	}

	// Detect type of the conversion job
	var rR renderRequest
	switch cj.RequestType {
	case "*service.imageRenderRequest":
		rR = &imageRenderRequest{}
	case "*service.pdfRenderRequest":
		rR = &pdfRenderRequest{}
	default:
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Error("invalid render target type, won't proceed")

		return nil
	}

	// Extract request details from the conversion job
	err = json.Unmarshal(cj.RequestData, &rR)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("error unmarshalling request data")

		return err
	}
	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Debug("extracted request data from conversion job")

	// Mark conversion job in 'processing' state and save the changes
	cj.markAsProcessing()
	err = cl.updateConversionJob(&cj)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("error: %v", err)

		return err
	}

	// Prepare conversion working directory, this is where we'll save the
	// resulting file before we upload it
	outputDir, err := ioutil.TempDir("", cj.Identifier)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("error: %v", err)

		return err
	}
	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Debug("prepared working directory for job")

	// Make sure we remove any generated files after we're done
	defer os.RemoveAll(outputDir)

	// Fulfill render request (perform actual conversion)
	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Info("start conversion process")
	outputLogs, outputFile, err := rR.fulfill(&cl, &cj, outputDir)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("error: %v", err)

		return err
	}
	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Info("completed conversion process")

	// Update conversion job with results
	cj.Logs = outputLogs
	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Debug("updated conversion job with logs")
	err = cl.updateConversionJob(&cj)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("error: %v", err)

		return err
	}

	// Upload the generated file to S3
	cj.StorageBucket = viper.GetString("worker.s3_bucket")
	cj.StorageKey = fmt.Sprintf("%s/%s", cj.Identifier, filepath.Base(outputFile))
	err = cl.storeFileS3(&cj, outputFile)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("error: %v", err)

		return err
	}

	// Update conversion job status and save the changes
	if err != nil {
		cj.markAsFailed()
	} else {
		cj.markAsSucceeded()
	}
	err = cl.updateConversionJob(&cj)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("error: %v", err)

		return err
	}

	return nil
}

// StartWorker starts the application background worker
func (c *Client) StartWorker() {
	concurrency := viper.GetSizeInBytes("worker.concurrency")
	maxRetries := viper.GetSizeInBytes("worker.max-retries")
	namespace := viper.GetString("redis.namespace")

	// Check for wkhtmltoimage installation
	_, version, erri := wkhtmltox.LookupConverter("wkhtmltoimage")
	log.Infof("using %s", version)
	if erri != nil {
		log.Errorln("unable to lookup wkhtmltoimage, make sure it's installed correctly")
	}

	// Check for wkhtmltopdf installation
	_, version, errp := wkhtmltox.LookupConverter("wkhtmltopdf")
	log.Infof("using %s", version)
	if erri != nil {
		log.Errorln("unable to lookup wkhtmltopdf, make sure it's installed correctly")
	}

	if erri != nil && errp != nil {
		log.Errorln("will not start workers due to errors")
	} else {
		log.Infof("concurrency set to %d", concurrency)
		log.Infof("maximum retries set to %d", maxRetries)
		pool := work.NewWorkerPool(workerContext{}, concurrency, namespace, c.redisPool)

		// Set job options
		maxFails := maxRetries + 1
		jobOptions := work.JobOptions{MaxFails: maxFails}

		// Assign jobs to queue
		log.Infof("registering '%s' queue", conversionQueue)
		pool.JobWithOptions(conversionQueue, jobOptions, (*workerContext).convert)

		// Start processing jobs
		log.Infof("waiting to pick up jobs placed on any registered queue")
		pool.Start()

		// Wait for a signal to quit
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, os.Kill)
		<-signalChan

		// Stop the pool
		pool.Stop()
	}
}
