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

package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

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
	}).Info("Picked up conversion job from queue")

	// Fetch all the job details
	cj, _, err := cl.fetchConversionJob(jid)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": jid,
		}).Errorf("Error: %v", err)

		// !!! //

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
		}).Error("Invalid render target type, won't proceed")

		// !!! //

		return nil
	}

	// Extract request details from the conversion job
	err = json.Unmarshal(cj.RequestData, &rR)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("Error unmarshalling request data")

		// !!! //

		return err
	}
	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Debug("Extracted request data from conversion job")

	// Mark conversion job in 'processing' state and save the changes
	cj.markAsProcessing()
	err = cl.updateConversionJob(&cj)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("Error: %v", err)

		// !!! //

		return err
	}

	// Prepare conversion working directory, this is where we'll save the
	// resulting file before we upload it
	outputDir, err := ioutil.TempDir("", cj.Identifier)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("Error: %v", err)

		// !!! //

		return err
	}
	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Debug("Prepared working directory for job")

	// Make sure we remove any generated files after we're done
	defer os.RemoveAll(outputDir)

	// Fulfill render request (perform actual conversion)
	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Info("Start conversion process")
	outputLogs, outputFile, err := rR.fulfill(&cl, &cj, outputDir)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("Error: %v", err)

		// !!! //

		return err
	}
	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Info("Completed conversion process")

	// Update conversion job with results
	cj.Logs = strings.TrimRight(string(outputLogs), "\r\n")
	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Debug("Updated conversion job with logs")
	err = cl.updateConversionJob(&cj)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("Error: %v", err)

		// !!! //

		return err
	}

	// Upload the generated file to S3
	cj.StorageRegion = viper.GetString("worker.s3_region")
	cj.StorageBucket = viper.GetString("worker.s3_bucket")
	cj.StorageKey = fmt.Sprintf("%s/%s", cj.Identifier, filepath.Base(outputFile))
	err = cl.storeFileS3(&cj, outputFile)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Errorf("Error: %v", err)

		// !!! //

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
		}).Errorf("Error: %v", err)

		// !!! //

		return err
	}

	return nil
}

// StartWorker starts the application background worker
func (c *Client) StartWorker() {
	concurrency := viper.GetSizeInBytes("worker.concurrency")
	namespace := viper.GetString("redis.namespace")

	// Check for wkhtmltoimage installation
	_, version, erri := wkhtmltox.LookupConverter("wkhtmltoimage")
	log.Infof("Using %s", version)
	if erri != nil {
		log.Errorln("Unable to lookup wkhtmltoimage, make sure it's installed correctly")
	}

	// Check for wkhtmltopdf installation
	_, version, errp := wkhtmltox.LookupConverter("wkhtmltopdf")
	log.Infof("Using %s", version)
	if erri != nil {
		log.Errorln("Unable to lookup wkhtmltopdf, make sure it's installed correctly")
	}

	if erri != nil && errp != nil {
		log.Errorln("Will not start workers due to errors")
	} else {
		log.Infof("Concurrency set to %d", concurrency)
		pool := work.NewWorkerPool(workerContext{}, concurrency, namespace, c.redisPool)

		// Assign queues to jobs
		log.Infof("Registering '%s' queue", conversionQueue)
		pool.Job(conversionQueue, (*workerContext).convert)

		// Start processing jobs
		log.Infof("Waiting to pick up jobs placed on any registered queue")
		pool.Start()

		// Wait for a signal to quit
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, os.Kill)
		<-signalChan

		// Stop the pool
		pool.Stop()
	}
}
