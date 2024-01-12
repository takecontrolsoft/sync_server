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

// Package sync configures and starts the sync server.
package main

import (
	"flag"
	"os"

	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/host"
	"github.com/takecontrolsoft/sync_server/server/services"
)

func main() {

	var port int
	var directory string

	portHelp := `TCP port number on witch the sync server can be reached. Defaults to 80.`
	flag.IntVar(&port, "p", 8080, portHelp)

	directoryHelp := `Storage path location for the synced files. It is required.
	This value should point to the directory where the uploaded files to be stored.
	Absolute path is required in DOS or UNC format.
	Make sure the server process has read/write access to this location.`
	flag.StringVar(&directory, "d", "", directoryHelp)

	flag.Parse()

	if argCount := len(os.Args[1:]); argCount == 0 {
		config.InitFromEnvVariables()
	} else {
		if directory == "" {
			logger.Fatal(config.ErrStoragePathEmpty)
		}
		config.PortNumber = port
		config.UploadDirectory = directory
	}

	logger.Info("Starting Sync server ...")

	logger.InfoF(" - port = %d\n", config.PortNumber)
	logger.InfoF(" - storage path = %s\n", config.UploadDirectory)
	services.Load()
	host.Run()
}
