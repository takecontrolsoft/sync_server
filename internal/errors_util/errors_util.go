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
package errors_util

import (
	"log"
	"os"

	"github.com/go-errors/errors"
)

func InitLogFile(logFileName string) {
	fLog, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(fLog)
	CrashOnError(err)
	defer saveLog(fLog)
}

func CrashOnError(err error) {
	if err != nil {
		log.Printf("ERROR CRASH: [%v]", err)
		panic(err)
	}
}

func LogError(err error) {
	if err != nil {
		log.Printf("ERROR CRASH: [%v]", err)
		println(err)
	}
}

func CrashWithError(err *errors.Error) {
	if err != nil {
		log.Printf("ERROR CRASH: [%v]", err)
		panic(err)
	}
}

func PrintOnError(err *errors.Error) {
	if err != nil {
		log.Printf("ERROR: [%v]", err)
		println(err)
	}
}

func RecoverOnError(err *errors.Error) {
	if err != nil {
		log.Printf("ERROR RECOVER: [%v]", err)
		recover()
	}
}

func saveLog(fLog *os.File) {
	fLog.Close()
}
