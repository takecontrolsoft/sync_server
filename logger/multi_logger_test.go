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

package logger_test

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/takecontrolsoft/go_multi_log/logger"
)

func init() {
}

func TestLogs(t *testing.T) {
	logger.Debug("Test log [debug] message")
	logger.Trace("Test log [trace] message")
	logger.Info("Test log [info] message")
	logger.Error(errors.Errorf("Test error.").Err)
	logger.Error("Test log [error] message")
	logger.Fatal("Test log [fatal] message")
}
