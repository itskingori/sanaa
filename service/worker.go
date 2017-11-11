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
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"

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

type context struct {
	client *Client
}

func (c *context) convert(job *work.Job) error {
	cl := NewClient()
	conn := cl.redisPool.Get()
	defer conn.Close()

	uuidStr := job.ArgString("uuid")
	cj, err := cl.getConversionJob(uuidStr)
	if err != nil {
		log.Error(err)
	}

	var rR renderRequest
	switch cj.RequestType {
	case "*service.imageRenderRequest":
		rR = &imageRenderRequest{}
	case "*service.pdfRenderRequest":
		rR = &pdfRenderRequest{}
	default:
		log.Errorf("Invalid target type: %s", cj.RequestType)
	}

	err = json.Unmarshal(cj.RequestData, &rR)
	if err != nil {
		return err
	}

	log.Infof("Starting processing of request %s", cj.Identifier)

	// Mark conversion job in processing state
	cj.StartedAt = time.Now().UTC().Format(time.RFC3339)
	cj.Status = "processing"

	// Save changes to conversion job: it's in 'processing' state at this point
	err = cl.updateConversionJob(&cj)
	if err != nil {
		log.Error(err)
	}

	// Prepare conversion working directory, this is where we'll save the
	// resulting file before we upload it
	outputDir, err := ioutil.TempDir("", cj.Identifier)
	if err != nil {
		log.Errorln(err)
	} else {
		log.Debugf("Prepared working directory for %s job", cj.Identifier)
	}
	defer os.RemoveAll(outputDir)

	// Fulfill render request (perform actual conversion)
	outputLogs, outputFile, gErr := rR.fulfill(&cl, &cj, outputDir)
	if gErr != nil {
		log.Errorf("Error fulfilling render request: %s", gErr)
	}

	// Update conversion job with the results and also update it's state to
	// reflect as such
	cj.Logs = strings.TrimRight(string(outputLogs), "\r\n")
	cj.EndedAt = time.Now().UTC().Format(time.RFC3339)
	if gErr != nil {
		cj.Status = "failed"
		log.Errorf("Failed to process request %s", cj.Identifier)
	} else {
		cj.OutputFile = outputFile
		cj.Status = "succeeded"
		log.Infof("Completed processing of request %s", cj.Identifier)
	}

	// Save changes to conversion job: it's either in 'failed' or 'succeeded'
	// state at this point
	err = cl.updateConversionJob(&cj)
	if err != nil {
		log.Error(err)
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
		pool := work.NewWorkerPool(context{}, concurrency, namespace, c.redisPool)

		// Assign queues to jobs
		log.Infof("Registering '%s' queue", conversionQueue)
		pool.Job(conversionQueue, (*context).convert)

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
