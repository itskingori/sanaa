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
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

const (
	// MinRequestTTL in the minimum TTL that we should allow to be set on requests
	MinRequestTTL = 300
)

type source struct {
	URL string `json:"url"`
}

type renderRequest interface {
	save(riq string, clt *Client) (ConversionJob, error)
	sourceURL() (*url.URL, error)
	fulfill(clt *Client, cj *ConversionJob, outputDir string) ([]byte, string, error)
}

type errorResponse struct {
	Identifier string `json:"uuid"`
	Message    string `json:"message"`
}

type renderResponse struct {
	Identifier string `json:"uuid"`
	CreatedAt  string `json:"created_at"`
	StartedAt  string `json:"started_at"`
	EndedAt    string `json:"ended_at"`
	ExpiresIn  int    `json:"expires_in"`
	FileURL    string `json:"file_url"`
	Status     string `json:"status"`
	Logs       string `json:"logs"`
}

func requestBadRequestResponse(w *http.ResponseWriter, r *http.Request, ers errorResponse) {
	log.WithFields(log.Fields{
		"uuid": ers.Identifier,
	}).Error(ers.Message)

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusBadRequest)

	encoder := json.NewEncoder((*w))
	encoder.SetEscapeHTML(false)
	encoder.Encode(&ers)

	log.WithFields(log.Fields{
		"uuid": ers.Identifier,
	}).Errorf("%d %s", http.StatusBadRequest, "Bad Request")
}

func requestInternalServerErrorResponse(w *http.ResponseWriter, r *http.Request, ers errorResponse) {
	log.WithFields(log.Fields{
		"uuid": ers.Identifier,
	}).Error(ers.Message)

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusInternalServerError)

	encoder := json.NewEncoder((*w))
	encoder.SetEscapeHTML(false)
	encoder.Encode(&ers)

	log.WithFields(log.Fields{
		"uuid": ers.Identifier,
	}).Errorf("%d %s", http.StatusInternalServerError, "Internal Server Error")
}

func requestNotFoundResponse(w *http.ResponseWriter, r *http.Request, ers errorResponse) {
	log.WithFields(log.Fields{
		"uuid": ers.Identifier,
	}).Error(ers.Message)

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusNotFound)

	encoder := json.NewEncoder((*w))
	encoder.SetEscapeHTML(false)
	encoder.Encode(&ers)

	log.WithFields(log.Fields{
		"uuid": ers.Identifier,
	}).Errorf("%d %s", http.StatusNotFound, "Not Found")
}

func requestCreatedResponse(w *http.ResponseWriter, r *http.Request, rrs renderResponse) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder((*w))
	encoder.SetEscapeHTML(false)
	encoder.Encode(&rrs)

	log.WithFields(log.Fields{
		"uuid": rrs.Identifier,
	}).Debugf("%d %s", http.StatusCreated, "Created")
}

func requestOKResponse(w *http.ResponseWriter, r *http.Request, rrs renderResponse) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)

	encoder := json.NewEncoder((*w))
	encoder.SetEscapeHTML(false)
	encoder.Encode(&rrs)

	log.WithFields(log.Fields{
		"uuid": rrs.Identifier,
	}).Debugf("%d %s", http.StatusOK, "OK")
}

func (clt *Client) renderHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ers errorResponse
		rrq renderRequest
	)

	params := mux.Vars(r)
	target := params["target"]
	rid := uuid.NewV4().String()

	switch target {
	case "image":
		rrq = &imageRenderRequest{}
	case "pdf":
		rrq = &pdfRenderRequest{}
	default:
		ers = errorResponse{
			Identifier: rid,
			Message:    fmt.Sprintf("Invalid %s render request", target),
		}
		requestBadRequestResponse(&w, r, ers)

		return
	}

	err := json.NewDecoder(r.Body).Decode(rrq)
	if err != nil {
		ers = errorResponse{
			Identifier: rid,
			Message:    fmt.Sprintf("Unable to unmarshal json to %s type", target),
		}
		requestBadRequestResponse(&w, r, ers)

		return
	}

	cj, err := rrq.save(rid, clt)
	if err != nil {
		ers = errorResponse{
			Identifier: rid,
			Message:    fmt.Sprintf("Unable to enqueue %s job", target),
		}
		requestInternalServerErrorResponse(&w, r, ers)

		return
	}
	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Infof("Enqueued render %s job", target)

	rrs, err := cj.generateRenderResponse(clt)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Error("Failed to generate render response")

		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Error(err)

		return
	}

	requestCreatedResponse(&w, r, rrs)
}

func (clt *Client) statusHandler(w http.ResponseWriter, r *http.Request) {
	var ers errorResponse

	params := mux.Vars(r)
	jid := params["uuid"]

	_, err := uuid.FromString(params["uuid"])
	if err != nil {
		ers = errorResponse{
			Identifier: jid,
			Message:    "Invalid job identifier",
		}
		requestBadRequestResponse(&w, r, ers)

		return
	}

	conn := clt.redisPool.Get()
	defer conn.Close()

	cj, found, err := clt.fetchConversionJob(jid)
	if err != nil {
		ers = errorResponse{
			Identifier: jid,
			Message:    "Unable to fetch conversion job",
		}
		requestInternalServerErrorResponse(&w, r, ers)

		return
	}

	if !found {
		ers = errorResponse{
			Identifier: jid,
			Message:    "Request not found on conversion queue",
		}
		requestNotFoundResponse(&w, r, ers)

		return
	}

	rrs, err := cj.generateRenderResponse(clt)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Error("Failed to generate render response")

		log.WithFields(log.Fields{
			"uuid": cj.Identifier,
		}).Error(err)

		return
	}

	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Info("Conversion job status check completed")

	requestOKResponse(&w, r, rrs)
}

func (cj *ConversionJob) generateRenderResponse(clt *Client) (renderResponse, error) {
	rrs := renderResponse{
		Identifier: cj.Identifier,
		CreatedAt:  cj.CreatedAt,
		StartedAt:  cj.StartedAt,
		EndedAt:    cj.EndedAt,
		ExpiresIn:  cj.ExpiresIn,
		Status:     cj.Status,
		Logs:       cj.Logs,
	}

	if cj.Status != "succeeded" {

		return rrs, nil
	}

	log.WithFields(log.Fields{
		"uuid": cj.Identifier,
	}).Debugln("Conversion job found completed")

	timeToExpire := 5 * time.Minute
	surl, err := clt.getFileS3SignedURL(cj, timeToExpire)
	if err != nil {

		return rrs, err
	}
	rrs.FileURL = surl

	return rrs, nil
}

// StartServer starts the application web server
func (clt *Client) StartServer() {
	requestTTL := viper.GetInt("server.request_ttl")
	log.Infof("Request TTL set to %d seconds", requestTTL)

	address := viper.GetString("server.binding_address")
	port := viper.GetInt("server.binding_port")
	binding := fmt.Sprintf("%s:%d", address, port)
	log.Infof("Listening on http://%s", binding)

	router := mux.NewRouter()
	router.HandleFunc("/render/{target}", clt.renderHandler).
		Headers("Content-Type", "application/json").
		Methods("POST")
	router.HandleFunc("/status/{uuid}", clt.statusHandler).
		Headers("Content-Type", "application/json").
		Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(binding, nil)
}
