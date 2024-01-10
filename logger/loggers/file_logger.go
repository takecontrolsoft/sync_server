/* Copyright 2024 Take Control - Software & Infrastructure

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

package loggers

import (
	"log"
	"os"
)

type FileLogger struct {
	LoggerInterface
}

func (logger *FileLogger) CrashOnError(log_err error) {
	fLog := setFileLog(logger)
	defer fLog.Close()
	log.Fatalf("ERROR: [%v]", log_err)
}

func (logger *FileLogger) LogError(log_err error) {
	fLog := setFileLog(logger)
	defer fLog.Close()
	log.Printf("ERROR: [%v]", log_err)
}

func (logger *FileLogger) LogMessage(message string) {
	fLog := setFileLog(logger)
	defer fLog.Close()
	log.Println(message)
}

func setFileLog(logger *FileLogger) *os.File {
	fLog, err := os.OpenFile("server_log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(fLog)
	return fLog
}
