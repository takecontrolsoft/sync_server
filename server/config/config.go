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
	"path/filepath"
	"runtime"

	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/go_multi_log/logger/levels"
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

// The name of the environment variable with the path for log files.
// This value should point to the directory where the log files to be stored.
// Absolute path is required in DOS or UNC format.
// Make sure the server process has read/write access to this location.
const LogPathVariable string = "LOG_PATH"

// The name of the environment variable for the log level.
// Allowed values are from 0 to 6.
// See package "go_multi_log": https://pkg.go.dev/github.com/takecontrolsoft/go_multi_log/logger/levels#LogLevel
const LogLevelVariable string = "LOG_LEVEL"

// Global variable for storage directory path
var UploadDirectory string

// Global variable for port number
var PortNumber int

// Global variable for log directory path
var LogPath string

// Global variable for log level
var LogLevel levels.LogLevel

// BinDirectory is the directory of the sync_server executable.
// Place exiftool and ffmpeg executables here so the server finds them.
var BinDirectory string

// InitBinDirectory sets BinDirectory to the executable's directory and
// prepends it to PATH so exiftool and ffmpeg are found when next to sync_server.
func InitBinDirectory() {
	exe, err := os.Executable()
	if err != nil {
		logger.Error(err)
		return
	}
	BinDirectory = filepath.Dir(exe)
	pathSep := string(os.PathListSeparator)
	if pathEnv := os.Getenv("PATH"); pathEnv != "" {
		_ = os.Setenv("PATH", BinDirectory+pathSep+pathEnv)
	} else {
		_ = os.Setenv("PATH", BinDirectory)
	}
	logger.InfoF("Bin directory (exiftool/ffmpeg): %s", BinDirectory)
}

// ExiftoolBinary returns the path to the exiftool executable (next to sync_server).
func ExiftoolBinary() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(BinDirectory, "exiftool.exe")
	}
	return filepath.Join(BinDirectory, "exiftool")
}

// FfmpegBinary returns the path to the ffmpeg executable (next to sync_server).
func FfmpegBinary() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(BinDirectory, "ffmpeg.exe")
	}
	return filepath.Join(BinDirectory, "ffmpeg")
}

// Initialize the variables [UploadDirectory], [PortNumber], [LogPath] and [LogLevel]
// from the environment variables [UploadPathVariable], [PortVariable], [LogPathVariable] and [LogLevelVariable].
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
	l, envSet := os.LookupEnv(LogPathVariable)
	if !envSet {
		l = ""
	}
	var n levels.LogLevel
	ll, envSet := os.LookupEnv(LogLevelVariable)
	if !envSet {
		n = levels.Info
	} else {
		_, err := fmt.Sscan(ll, &n)
		if err != nil {
			logger.Fatal(err)
		}
	}
	UploadDirectory = d
	PortNumber = p
	LogPath = l
	LogLevel = n
	logger.InfoF("Server port: %d", PortNumber)
	logger.InfoF(fmt.Sprintf("Storage path: %s", UploadDirectory))
	logger.InfoF(fmt.Sprintf("Log path: %s", LogPath))
	logger.InfoF(fmt.Sprintf("Log level: %s", LogLevel.String()))

}
