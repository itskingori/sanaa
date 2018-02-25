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
	"net/url"
	"path/filepath"

	"github.com/itskingori/go-wkhtml/wkhtmltox"

	log "github.com/sirupsen/logrus"
)

type pdfRenderRequest struct {
	Source source               `json:"source"`
	Target wkhtmltox.PDFOptions `json:"target"`
}

func (rr *pdfRenderRequest) save(riq string, c *Client) (ConversionJob, error) {
	cj, err := c.createAndSaveConversionJob(riq, rr)
	if err != nil {
		log.Error(err)

		return cj, err
	}

	return cj, nil
}

func (rr *pdfRenderRequest) sourceURL() (*url.URL, error) {
	u, err := url.Parse(rr.Source.URL)
	if err != nil {
		log.Error(err)

		return u, err
	}

	return u, nil
}

func (rr *pdfRenderRequest) fulfill(c *Client, cj *ConversionJob, outputDir string) ([]byte, string, error) {
	var (
		outputFile string
		outputLogs []byte
	)

	opts := rr.Target
	pfs := wkhtmltox.NewPDFFlagSetFromOptions(&opts)
	outputFile = filepath.Join(outputDir, "file.pdf")
	inputURL := rr.Source.URL
	outputLogs, err := pfs.Generate(inputURL, outputFile)
	if err != nil {
		log.Error(err)

		return outputLogs, outputFile, err
	}

	return outputLogs, outputFile, nil
}
