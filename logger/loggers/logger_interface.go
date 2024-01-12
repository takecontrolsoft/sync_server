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
)

type LoggerInterface interface {
	Log(level int, arg any)
	LogF(level int, format string, args ...interface{})
}

const (
	DebugLevel   = 1
	TraceLevel   = 2
	InfoLevel    = 3
	WarningLevel = 4
	ErrorLevel   = 5
	FatalLevel   = 6
)

func multi_log(level int, arg any) {
	format := fmt.Sprintf("%s: [%s]", logLevelName(level), "%v")
	multi_logF(level, format, arg)
}

func multi_logF(level int, format string, args ...interface{}) {
	log.Printf(format, args...)
}

func logLevelName(level int) string {
	switch logLevel := level; logLevel {
	case DebugLevel:
		return "DEBUG"
	case TraceLevel:
		return "TRACE"
	case InfoLevel:
		return "INFO"
	case WarningLevel:
		return "WARNING"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
