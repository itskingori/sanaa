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
	"net/url"
	"time"

	"github.com/itskingori/go-wkhtml/pdf"

	log "github.com/sirupsen/logrus"
)

type pdfRenderRequest struct {
	Source source      `json:"source"`
	Target pdf.Options `json:"target"`
}

func (r *pdfRenderRequest) save(c *Client) (ConversionJob, error) {
	cj, err := c.saveRequestJobDetails(r)
	if err != nil {
		log.Error(err)
	}

	return cj, nil
}

func (r *pdfRenderRequest) sourceURL() (*url.URL, error) {
	u, err := url.Parse(r.Source.URL)
	if err != nil {
		log.Error(err)
	}

	return u, nil
}

func (r *pdfRenderRequest) runConversion(c *Client, cj *ConversionJob) error {
	cj.StartedAt = time.Now().String()
	log.Infof("Starting processing of request %s", cj.Identifier)
	err := c.updateRequestJobDetails(cj)
	if err != nil {
		log.Error(err)
	}

	// COVERT HERE!

	cj.EndedAt = time.Now().String()
	err = c.updateRequestJobDetails(cj)
	if err != nil {
		log.Error(err)
	}
	log.Infof("Completed processing of request %s", cj.Identifier)

	return nil
}
