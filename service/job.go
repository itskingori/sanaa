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
	"reflect"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/work"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// ConversionJob is a mapping of a conversion job's attributes
type ConversionJob struct {
	Identifier  string `redis:"uuid"`
	CreatedAt   string `redis:"created_at"`
	StartedAt   string `redis:"started_at"`
	EndedAt     string `redis:"ended_at"`
	ExpiresIn   int    `redis:"expires_in"`
	Status      string `redis:"status"`
	Logs        string `redis:"logs"`
	OutputFile  string `redis:"output_file"`
	RequestType string `redis:"request_type"`
	RequestData []byte `redis:"request_data"`
}

func generateJobID(name string) uuid.UUID {
	// Now we create our uuid using the random uuid and a namespace value which
	// will be determined from the domain part of the source URL.
	//
	// This allow us to generate several UUIDs with a low probability of
	// collision. See: https://tools.ietf.org/html/rfc4122#section-4.3
	return uuid.NewV5(uuid.NewV4(), name)
}

func jobKey(u uuid.UUID) string {
	return fmt.Sprintf("%s:request:%s", viper.GetString("redis.namespace"), u)
}

func (c *Client) enqueueConversionJob(u uuid.UUID) error {
	_, err := c.enqueuer.Enqueue(conversionQueue, work.Q{"uuid": u})
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (c *Client) createAndSaveConversionJob(rR renderRequest) (ConversionJob, error) {
	cj := ConversionJob{}
	rt := viper.GetInt("server.request_ttl")

	su, err := rR.sourceURL()
	if err != nil {
		return cj, err
	}

	uid := generateJobID(su.Host)
	key := jobKey(uid)

	serializedRequest, err := json.Marshal(rR)
	if err != nil {
		return cj, err
	}

	cj.Identifier = uid.String()
	cj.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	cj.ExpiresIn = rt
	cj.Status = "pending"
	cj.RequestType = reflect.TypeOf(rR).String()
	cj.RequestData = serializedRequest

	conn := c.redisPool.Get()
	defer conn.Close()

	_, err = conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(&cj)...)
	if err != nil {
		return cj, err
	}

	_, err = conn.Do("EXPIRE", key, rt)
	if err != nil {
		return cj, err
	}

	err = c.enqueueConversionJob(uid)
	if err != nil {
		return cj, err
	}

	return cj, nil
}

func (c *Client) getConversionJob(uuidStr string) (ConversionJob, error) {
	var cj ConversionJob

	conn := c.redisPool.Get()
	defer conn.Close()

	uid, err := uuid.FromString(uuidStr)
	if err != nil {
		return cj, err
	}

	v, err := redis.Values(conn.Do("HGETALL", jobKey(uid)))
	if err != nil {
		return cj, err
	}

	err = redis.ScanStruct(v, &cj)
	if err != nil {
		return cj, err
	}

	return cj, err
}

func (c *Client) updateConversionJob(cj *ConversionJob) error {
	conn := c.redisPool.Get()
	defer conn.Close()

	uid, err := uuid.FromString(cj.Identifier)
	if err != nil {
		return err
	}

	job := *cj
	key := jobKey(uid)
	_, err = conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(&job)...)
	if err != nil {
		log.Errorf("Error saving changes to %s job", cj.Identifier)

		return err
	}

	log.Debugf("Saved changes to %s job", cj.Identifier)

	return nil
}
