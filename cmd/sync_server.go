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

// Package sync_server configures and starts the sync server.
package sync_server

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/takecontrolsoft/go_multi_log/logger/levels"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/host"
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

		scanner := bufio.NewScanner(os.Stdin)
		if directory == "" {
			fmt.Printf("Please enter storage path location in DOS or UNC format:")
			for scanner.Scan() {
				directory = scanner.Text()
				if directory != "" {
					if _, err := os.Stat(directory); os.IsNotExist(err) {
						fmt.Printf("Directory '%s' does not exist. Please enter a valid path.", directory)
					} else {
						break
					}
				}
				fmt.Println("")
			}
		}
		config.UploadDirectory = directory

		if port == 80 {
			fmt.Printf("Please enter TCP port number. (default: 80):")
			for scanner.Scan() {
				v := scanner.Text()
				if v != "" {
					_, err := fmt.Sscan(v, &port)
					if err != nil {
						fmt.Println(v, err, reflect.TypeOf(port))
					} else {
						break
					}
				} else {
					port = 80
					break
				}
			}
		}
		config.PortNumber = port

		if logPath == "" {
			fmt.Printf("Please enter log files location:")
			for scanner.Scan() {
				v := scanner.Text()
				if v != "" {
					_, err := fmt.Sscan(v, &logPath)
					if err != nil {
						fmt.Println(v, err, reflect.TypeOf(logPath))
					} else {
						break
					}
				} else {
					logPath = ""
					break
				}
			}
		}
		config.LogPath = logPath

		if logLevel == 3 {
			fmt.Printf("Please enter log level:")
			for scanner.Scan() {
				n := scanner.Text()
				if n != "" {
					_, err := fmt.Sscan(n, &logLevel)
					if err != nil {
						fmt.Println(n, err, reflect.TypeOf(logLevel))
					} else {
						break
					}
				} else {
					logLevel = 3
					break
				}
			}
		}
		config.LogLevel = levels.LogLevel(logLevel)
	}
	fmt.Println("Starting Sync server ...")

	fmt.Printf(" - port = %d\n", config.PortNumber)
	fmt.Printf(" - storage path = %s\n", config.UploadDirectory)
	fmt.Printf(" - log files path = %s\n", config.LogPath)
	fmt.Printf(" - log level = %s\n", config.LogLevel.String())
	host.Run()
}
