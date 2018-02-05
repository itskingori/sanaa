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
	Identifier    string `redis:"uuid"`
	CreatedAt     string `redis:"created_at"`
	StartedAt     string `redis:"started_at"`
	EndedAt       string `redis:"ended_at"`
	ExpiresIn     int    `redis:"expires_in"`
	Status        string `redis:"status"`
	Logs          string `redis:"logs"`
	StorageRegion string `redis:"storage_region"`
	StorageBucket string `redis:"storage_bucket"`
	StorageKey    string `redis:"storage_key"`
	RequestType   string `redis:"request_type"`
	RequestData   []byte `redis:"request_data"`
}

func generateJobKey(jid string) string {
	key := fmt.Sprintf("%s:request:%s", viper.GetString("redis.namespace"), jid)

	return key
}

func (cj *ConversionJob) markAsProcessing() {
	cj.StartedAt = time.Now().UTC().Format(time.RFC3339)
	cj.Status = "processing"

	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Info("Marked conversion job as 'processing'")
}

func (cj *ConversionJob) markAsFailed() {
	cj.EndedAt = time.Now().UTC().Format(time.RFC3339)
	cj.Status = "failed"

	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Info("Marked conversion job as 'failed'")
}

func (cj *ConversionJob) markAsSucceeded() {
	cj.EndedAt = time.Now().UTC().Format(time.RFC3339)
	cj.Status = "succeeded"

	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Info("Marked conversion job as 'succeeded'")
}

func (clt *Client) enqueueConversionJob(riq string) error {
	_, err := clt.enqueuer.Enqueue(conversionQueue, work.Q{"uuid": riq})
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (clt *Client) createAndSaveConversionJob(rid string, rR renderRequest) (ConversionJob, error) {
	cj := ConversionJob{}
	rt := viper.GetInt("server.request_ttl")
	key := generateJobKey(rid)

	serializedRequest, err := json.Marshal(rR)
	if err != nil {
		return cj, err
	}

	cj.Identifier = rid
	cj.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	cj.ExpiresIn = rt
	cj.Status = "pending"
	cj.RequestType = reflect.TypeOf(rR).String()
	cj.RequestData = serializedRequest

	conn := clt.redisPool.Get()
	defer conn.Close()

	conn.Send("HMSET", redis.Args{}.Add(key).AddFlat(&cj)...)
	conn.Send("EXPIRE", key, rt)
	conn.Flush()

	_, err = conn.Receive()
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": rid,
		}).Error("Error saving conversion job")

		return cj, err
	}

	_, err = conn.Receive()
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": rid,
		}).Error("Error setting conversion job expiry")

		return cj, err
	}

	err = clt.enqueueConversionJob(rid)
	if err != nil {
		return cj, err
	}

	return cj, nil
}

func (clt *Client) fetchConversionJob(jid string) (ConversionJob, bool, error) {
	conn := clt.redisPool.Get()
	defer conn.Close()

	cj := ConversionJob{}
	found := false

	value, err := redis.Values(conn.Do("HGETALL", generateJobKey(jid)))
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": jid,
		}).Error("Unable to fetch values from redis")

		return cj, found, err
	}

	if len(value) == 0 {

		return cj, found, err
	}

	err = redis.ScanStruct(value, &cj)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": jid,
		}).Error("Unable to unmarshall values to job")

		return cj, found, err
	}
	found = true

	log.WithFields(log.Fields{
		"uuid": jid,
	}).Debug("Fetched conversion job details")

	return cj, found, err
}

func (clt *Client) updateConversionJob(cj *ConversionJob) error {
	conn := clt.redisPool.Get()
	defer conn.Close()

	uid, err := uuid.FromString(cj.Identifier)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Error("Unable to parse job identifier")

		return err
	}

	job := *cj
	key := generateJobKey(uid.String())
	_, err = conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(&job)...)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Error("Error saving conversion job changes")

		return err
	}

	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Debug("Saved conversion job changes")

	return nil
}
