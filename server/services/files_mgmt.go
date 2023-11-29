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

// Package services registers the APIs for the sync server.
package services

import (
	"fmt"
	"net/http"

	"github.com/takecontrolsoft/sync_server/server/config"
	host "github.com/takecontrolsoft/sync_server/server/host"
	"github.com/takecontrolsoft/sync_server/server/impl"
)

type FilesManagementService struct{}

func (s FilesManagementService) Host() bool {
	fmt.Println("FilesManagementService::Host()")
	http.HandleFunc("/upload", impl.UploadHandler)

	fs := http.FileServer(http.Dir(config.UploadDirectory))
	http.Handle("/files", http.StripPrefix("/files", fs))

	return true
}

func init() {
	host.RegisterWebService(FilesManagementService{})
}
