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
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/itskingori/go-wkhtml/wkhtmltox"

	log "github.com/sirupsen/logrus"
)

type imageRenderRequest struct {
	Source source                 `json:"source"`
	Target wkhtmltox.ImageOptions `json:"target"`
}

func (rr *imageRenderRequest) save(riq string, c *Client) (ConversionJob, error) {
	cj, err := c.createAndSaveConversionJob(riq, rr)
	if err != nil {
		log.Error(err)

		return cj, err
	}

	return cj, nil
}

func (rr *imageRenderRequest) sourceURL() (*url.URL, error) {
	u, err := url.Parse(rr.Source.URL)
	if err != nil {
		log.Error(err)

		return u, err
	}

	return u, nil
}

func (rr *imageRenderRequest) fulfill(c *Client, cj *ConversionJob, outputDir string) ([]byte, string, error) {
	var (
		outputFile string
		outputLogs []byte
	)

	opts := rr.Target
	ifs := wkhtmltox.NewImageFlagSetFromOptions(&opts)
	format, _ := ifs.GetFormat()
	outputFile = filepath.Join(outputDir, fmt.Sprintf("file.%s", format))
	inputURL := rr.Source.URL
	outputLogs, err := ifs.Generate(inputURL, outputFile)
	if err != nil {
		log.Error(err)

		return outputLogs, outputFile, err
	}

	return outputLogs, outputFile, nil
}
