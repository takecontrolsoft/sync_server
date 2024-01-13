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
	"strings"
)

const (
	DebugLevel   int = 1
	TraceLevel   int = 2
	InfoLevel    int = 3
	WarningLevel int = 4
	ErrorLevel   int = 5
	FatalLevel   int = 6
)

type LoggerInterface interface {
	Log(level int, arg any)
	LogF(level int, format string, args ...interface{})
	GetLevel()
	SetLevel(level int)
}

type loggerType struct {
	LoggerInterface
	Level  int
	Format string
}

func (logger *loggerType) IsLogLevelAllowed(level int) bool {
	return level >= logger.Level
}

func (logger *loggerType) GetLevel() int {
	return logger.Level
}

func (logger *loggerType) SetLevel(level int) {
	logger.Level = level
}

func (logger *loggerType) multi_log(level int, arg any) {
	format := fmt.Sprintf("%s: [%s]", strings.ToUpper(GetLogLevelName(level)), "%v")
	logger.multi_logF(level, format, arg)
}

func (logger *loggerType) multi_logF(level int, format string, args ...interface{}) {
	log.Printf(format, args...)
}

func GetLogLevelName(level int) string {
	switch logLevel := level; logLevel {
	case DebugLevel:
		return "Debug"
	case TraceLevel:
		return "Trace"
	case InfoLevel:
		return "Info"
	case WarningLevel:
		return "Warning"
	case ErrorLevel:
		return "Error"
	case FatalLevel:
		return "Fatal"
	default:
		return "Unknown"
	}
}
