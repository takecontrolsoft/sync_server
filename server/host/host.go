/* Copyright 2023 Take Control - Software & Infrastructure

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package host provides functions for registering the web service APIs
// and running the sync server.
package host

import (
	"fmt"
	"net/http"

	"github.com/takecontrolsoft/go_multi_log/logger"

	"github.com/takecontrolsoft/sync_server/server/config"
)

// Service interface that should be implemented by all the services.
type WebService interface {

	// Attach http handles to http routs.
	Host() bool
}

// Web services that to be hosted.
var webServices = []WebService{}

// Register a web service type.
func RegisterWebService(s WebService) {
	webServices = append(webServices, s)
}

// Start sync server and host all the registered web services.
func Run() {
	for _, service := range webServices {
		service.Host()
	}
	logger.Info("Sync server started.")
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.PortNumber), nil)
	if err != nil {
		logger.Fatal(err)
	}

}
