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
	"github.com/takecontrolsoft/go_multi_log/logger/levels"
	"github.com/takecontrolsoft/go_multi_log/logger/loggers"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/host"
	"github.com/takecontrolsoft/sync_server/server/services"
)

func main() {

	var port int
	var directory string
	var logPath string
	var logLevel int

	portHelp := `TCP port number on witch the sync server can be reached. Defaults to 80.`
	flag.IntVar(&port, "p", 8080, portHelp)

	directoryHelp := `Storage path location for the synced files. It is required.
	This value should point to the directory where the uploaded files to be stored.
	Absolute path is required in DOS or UNC format.
	Make sure the server process has read/write access to this location.`
	flag.StringVar(&directory, "d", "", directoryHelp)

	logPathHelp := `Path location for the log files. 
	If not set, the log files will be stored to the executable file location.
	Absolute path is required in DOS or UNC format.
	Make sure the server process has read/write access to this location.`
	flag.StringVar(&logPath, "l", "", logPathHelp)

	logLevelHelp := `Log level. 
	If not set, the log level will be set to Info by default.
	Allowed values are from 0 to 6.
    See package "go_multi_log": https://pkg.go.dev/github.com/takecontrolsoft/go_multi_log/logger/levels#LogLevel`
	flag.IntVar(&logLevel, "n", 3, logLevelHelp)

	flag.Parse()

	if argCount := len(os.Args[1:]); argCount == 0 {
		config.InitFromEnvVariables()
	} else {
		if directory == "" {
			logger.Fatal(config.ErrStoragePathEmpty)
		}
		config.PortNumber = port
		config.UploadDirectory = directory
		config.LogPath = logPath
		config.LogLevel = levels.LogLevel(logLevel)
	}

	level := config.LogLevel
	fileOptions := loggers.FileOptions{
		Directory:     config.LogPath,
		FilePrefix:    "sync_server",
		FileExtension: ".log",
	}

	f := loggers.NewFileLogger(level, "", fileOptions)
	err := logger.RegisterLogger("log_key", f)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Starting Sync server ...")

	config.InitBinDirectory()
	logger.InfoF(" - port = %d\n", config.PortNumber)
	logger.InfoF(" - storage path = %s\n", config.UploadDirectory)
	logger.InfoF(" - log path = %s\n", config.LogPath)
	logger.InfoF(" - log level = %s\n", config.LogLevel.String())
	services.Load()
	host.Run()
}
