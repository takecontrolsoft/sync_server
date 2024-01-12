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

// Package config provides functions for setting
// the initial values of server parameters.
package config

import (
	"fmt"
	"os"

	"github.com/takecontrolsoft/go_multi_log/logger"
)

// The maximum stream size that is allowed to be uploaded to the server.
// The size is set to maximum 5GB.
const MaxUploadFileSize int64 = 5 * 1024 * 1024 * 1024

// The name of the environment variable with the storage path.
// This value should point to the directory where the uploaded files to be stored.
// Absolute path is required in DOS or UNC format.
// Make sure the server process has read/write access to this location.
const UploadPathVariable string = "SYNC_STORAGE_PATH"

// The name of the environment variable with TCP port number
// on witch the server can be reached.
const PortVariable = "SYNC_SERVER_PORT"

// Global variable for storage directory path
var UploadDirectory string

// Global variable for port number
var PortNumber int

// Initialize the variables [UploadDirectory] and [PortNumber]
// from the environment variables [UploadPathVariable] and [PortVariable].
func InitFromEnvVariables() {
	d, envSet := os.LookupEnv(UploadPathVariable)
	if !envSet {
		logger.Fatal(ErrEnvVariableNotSet(UploadPathVariable))
	}
	if d == "" {
		logger.Fatal(ErrEnvVariableSetEmpty(UploadPathVariable))
	}
	port, envSet := os.LookupEnv(PortVariable)
	if !envSet {
		logger.Fatal(ErrEnvVariableNotSet(PortVariable))
	}
	var p int
	if port == "" {
		logger.Fatal(ErrEnvVariableSetEmpty(PortVariable))
	} else {
		_, err := fmt.Sscan(port, &p)
		if err != nil {
			logger.Fatal(err)
		}
	}

	UploadDirectory = d
	PortNumber = p
	logger.InfoF("Server port: %d", PortNumber)
	logger.InfoF(fmt.Sprintf("Storage path: %s", UploadDirectory))
}
