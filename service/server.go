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
	save(c *Client) (uuid.UUID, error)
	sourceURL() (*url.URL, error)
	runConversion(c *Client, cj *ConversionJob) error
}

func (c *Client) renderHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	target := params["target"]

	var rR renderRequest
	switch target {
	case "image":
		rR = &imageRenderRequest{}
	case "pdf":
		rR = &pdfRenderRequest{}
	default:
		log.Errorf("Invalid target: %s", target)
	}

	err := json.NewDecoder(r.Body).Decode(rR)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Received render %s request", target)

	uid, err := rR.save(c)
	if err != nil {
		log.Error(err)
	}
	log.Infof("Enqueued render %s job, uuid: %s", target, uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rR)
}

// StartServer starts the application web server
func (c *Client) StartServer() {
	requestTTL := viper.GetInt("server.request_ttl")
	log.Infof("Request TTL set to %d seconds", requestTTL)

	address := viper.GetString("server.binding_address")
	port := viper.GetInt("server.binding_port")
	binding := fmt.Sprintf("%s:%d", address, port)
	log.Infof("Listening on http://%s", binding)

	router := mux.NewRouter()
	router.HandleFunc("/render/{target}", c.renderHandler).
		Headers("Content-Type", "application/json").
		Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(binding, nil)
}
