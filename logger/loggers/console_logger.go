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
	LoggerInterface
}

func (logger *ConsoleLogger) CrashOnError(err error) {
	if err != nil {
		log.SetOutput(os.Stdout)
		log.Fatalf("ERROR: [%v]\n", err)
	}
}

func (logger *ConsoleLogger) LogError(err error) {
	if err != nil {
		log.SetOutput(os.Stdout)
		log.Printf("ERROR: [%v]\n", err)
	}
}

func (logger *ConsoleLogger) LogMessage(message string) {
	log.SetOutput(os.Stdout)
	log.Println(message)
}
