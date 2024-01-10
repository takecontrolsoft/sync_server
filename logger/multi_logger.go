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

// Package errors_util provides functions for errors and logs.
package logger

import (
	"sync"

	"github.com/takecontrolsoft/logger/loggers"
)

var lock = &sync.Mutex{}

type MultiLog struct {
	registered_loggers []loggers.LoggerInterface
}

var mLogger *MultiLog

func getLogger() *MultiLog {
	if mLogger == nil {
		lock.Lock()
		defer lock.Unlock()
		if mLogger == nil {
			mLogger = &MultiLog{
				registered_loggers: []loggers.LoggerInterface{
					&loggers.ConsoleLogger{},
					&loggers.FileLogger{}},
			}
		}
	}

	return mLogger
}

func CrashOnError(log_err error) {
	if log_err != nil {
		mLogger = getLogger()
		for _, logger := range mLogger.registered_loggers {
			logger.CrashOnError(log_err)
		}
	}
}

func LogError(log_err error) {
	if log_err != nil {
		mLogger = getLogger()
		for _, logger := range mLogger.registered_loggers {
			logger.LogError(log_err)
		}
	}
}

func LogMessage(message string) {
	mLogger = getLogger()
	for _, logger := range mLogger.registered_loggers {
		logger.LogMessage(message)
	}
}
