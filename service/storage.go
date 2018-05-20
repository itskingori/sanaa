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
	"bytes"
	"context"
	"time"

	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"

	log "github.com/sirupsen/logrus"
)

func (cl *Client) storeFileS3(cj *ConversionJob, filePath string) error {
	svc := s3.New(cl.awsSession)

	// Create a context with a timeout that will abort the upload if it takes more
	// than the passed in timeout
	timeout := 60 * time.Second
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, timeout)

	// Ensure the context is canceled to prevent leaking. See context package for
	// more information, https://golang.org/pkg/context/
	defer cancelFn()

	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Debug("read file from working directory")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Info("start upload of file to S3")

	// Uploads the object to S3 ... the Context will interrupt the request if the
	// timeout expires
	_, err = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(cj.StorageBucket),
		Key:    aws.String(cj.StorageKey),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned
			log.WithFields(log.Fields{
				"uuid": cj.Identifier,
			}).Errorf("upload cancelled due to timeout: %v", err)
		} else {
			log.WithFields(log.Fields{
				"uuid": cj.Identifier,
			}).Errorf("failed to upload file: %v", err)
		}

		return err
	}

	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Info("completed upload of file to S3")

	return nil
}

func (cl *Client) getFileS3SignedURL(cj *ConversionJob, exp time.Duration) (string, error) {
	svc := s3.New(cl.awsSession)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &cj.StorageBucket,
		Key:    &cj.StorageKey,
	})

	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Debugln("generating pre-signed url to rendered file")

	url, err := req.Presign(exp)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Error("failed to pre-sign url")

		return url, err
	}

	return url, nil
}
