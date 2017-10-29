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
	"os"
	"os/signal"

	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/work"
	"github.com/satori/go.uuid"
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
	uid, err := uuid.FromString(job.ArgString("uuid"))
	if err != nil {
		return err
	}

	cl := NewClient()
	conn := cl.redisPool.Get()
	defer conn.Close()

	v, err := redis.Values(conn.Do("HGETALL", jobKey(uid)))
	if err != nil {
		return err
	}

	var cj ConversionJob
	err = redis.ScanStruct(v, &cj)
	if err != nil {
		return err
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

	err = rR.runConversion(&cl, &cj)
	if err != nil {
		return err
	}

	return nil
}

// StartWorker starts the application background worker
func (c *Client) StartWorker() {
	concurrency := viper.GetSizeInBytes("worker.concurrency")
	namespace := viper.GetString("redis.namespace")

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
