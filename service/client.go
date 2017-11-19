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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/work"
	"github.com/spf13/viper"
)

// Client is the application client
type Client struct {
	awsSession *session.Session
	enqueuer   *work.Enqueuer
	redisPool  *redis.Pool
}

// NewClient creates an initialized application client
func NewClient() Client {
	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			host := viper.GetString("redis.host")
			port := viper.GetInt("redis.port")

			return redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		},
	}
	enqueuer := work.NewEnqueuer(viper.GetString("redis.namespace"), redisPool)

	region := viper.GetString("worker.s3_region")
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	return Client{
		awsSession: sess,
		enqueuer:   enqueuer,
		redisPool:  redisPool,
	}
}
