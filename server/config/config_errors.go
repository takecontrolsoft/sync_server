/*
	Copyright 2023 Take Control - Software & Infrastructure

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

package config

import "github.com/go-errors/errors"

// An error for not set value for an environment variable that is needed to run the sync server.
func ErrEnvVariableNotSet(envVariableName string) *errors.Error {
	return errors.Errorf("Environment variable %s is not set.", envVariableName)
}

// An error for empty value for an environment variable that is needed to run the sync server.
func ErrEnvVariableSetEmpty(envVariableName string) *errors.Error {
	return errors.Errorf("Environment variable %s is empty.", envVariableName)
}

// An error for empty storage path.
var ErrStoragePathEmpty = errors.Errorf("Storage path is empty.")
