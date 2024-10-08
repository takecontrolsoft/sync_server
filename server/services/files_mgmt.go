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
	"net/http"

	"github.com/takecontrolsoft/go_multi_log/logger"
	host "github.com/takecontrolsoft/sync_server/server/host"
	"github.com/takecontrolsoft/sync_server/server/impl"
)

type FilesManagementService struct {
	host.WebService
}

func (s FilesManagementService) Host() bool {
	logger.Info("FilesManagementService hosted")

	http.HandleFunc("/upload", impl.UploadHandler)

	//fs := http.FileServer(http.Dir(config.UploadDirectory))
	//http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/folders", impl.GetFoldersHandler)

	http.HandleFunc("/files", impl.GetFilesHandler)

	http.HandleFunc("/delete-all", impl.DeleteAllHandler)

	return true
}

func init() {
	host.RegisterWebService(FilesManagementService{})
}
