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
	"fmt"
	"log"
	"os"
)

type FileLogger struct {
	LoggerInterface
	FilePrefix string
}

func (logger *FileLogger) Log(level int, arg any) {
	fLog := setFileLog(logger)
	defer fLog.Close()
	multi_log(level, arg)
}

func (logger *FileLogger) LogF(level int, format string, args ...interface{}) {
	fLog := setFileLog(logger)
	defer fLog.Close()
	multi_logF(level, format, args...)
}

func setFileLog(logger *FileLogger) *os.File {
	fLog, err := os.OpenFile(fmt.Sprintf("%s.log", logger.FilePrefix), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(fLog)
	return fLog
}
