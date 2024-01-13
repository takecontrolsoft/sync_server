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

type ConsoleLogger struct {
	loggerType
}

func NewConsoleLoggerDefault() *ConsoleLogger {
	return &ConsoleLogger{
		loggerType: loggerType{Level: InfoLevel},
	}
}

func NewConsoleLogger(level int, format string) *ConsoleLogger {
	return &ConsoleLogger{
		loggerType: loggerType{Level: level, Format: format},
	}
}

func (logger *ConsoleLogger) Log(level int, arg any) {
	if logger.IsLogLevelAllowed(level) {
		log.SetOutput(os.Stdout)
		logger.multi_log(level, arg)
	}
}

func (logger *ConsoleLogger) LogF(level int, format string, args ...interface{}) {
	if logger.IsLogLevelAllowed(level) {
		log.SetOutput(os.Stdout)
		logger.multi_logF(level, format, args...)
	}
}